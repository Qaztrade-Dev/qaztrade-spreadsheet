class HeaderCell {
  constructor(key, range, values, groupKey = '') {
    this.key = key;
    this.range = range;
    this.values = values;
    this.groupKey = groupKey;
  }
}

const parents = {
  '№': {
    parent: 'root',
    parentKey: 'root',
  },
  'Дистрибьюторский договор': {
    parent: '№',
    parentKey: '№',
  },
  'контракт на поставку': {
    parent: 'Дистрибьюторский договор',
    parentKey: 'Дистрибьюторский договор.№',
  },
};

const headerCells = getHeaderCells();

function my() {
  childKey = 'Дистрибьюторский договор'

  a = getParentValues('Дистрибьюторский договор')
  Logger.log(a)
}

function getParentValues(childKey) {
  const parentHeaderCell = getHeaderCell(parents[childKey].parentKey);
  const values = getLevelValues(parentHeaderCell);
  
  const filteredValues = values.reduce((nonEmptyStrings, str) => {
    if (typeof str === 'string' && str.trim() !== '') {
      nonEmptyStrings.push(str);
    }
    return nonEmptyStrings;
  }, []);

  return filteredValues;
}

function getHeaderCell(cellKey) {
  const keys = cellKey.split('.');
  let cell = headerCells[keys[0]];

  for (let i = 1; i < keys.length; i++) {
    const key = keys[i];
    cell = cell.values[key];
  }

  return cell;
}

function getHeaderCells() {
  const sheet = SpreadsheetApp.getActiveSpreadsheet();
  const headerValues = sheet.getRangeByName('Header').getValues();

  const topLevel = headerValues[0];
  const lowLevel = headerValues[1];

  const cellMap = {};
  for (let i = 0; i < topLevel.length; i++) {
    if (topLevel[i] === '') {
      continue;
    }

    const values = {};
    let r = i;

    for (let j = i; j < lowLevel.length; j++) {
      if (lowLevel[j] == '') {
        break;
      }

      if (!(topLevel[j] == '' || i == j)) {
        break;
      }

      values[lowLevel[j]] = new HeaderCell(lowLevel[j], [j + 1, j + 1], {}, topLevel[i]);
      r = j;
    }

    cellMap[topLevel[i]] = new HeaderCell(topLevel[i], [i + 1, r + 1], values, topLevel[i]);
  }

  return cellMap;
}

function getLevelValues(headerCell) {
  const upperBound = 6;
  const columnNum = headerCell.range[0];

  const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheets()[0];
  const rows = sheet.getRange(upperBound, columnNum, sheet.getMaxRows(), 1).getValues();

  const values = [];
  for (let i = 0; i < rows.length; i++) {
    const value = rows[i][0].toString();
    if (value === '') {
      continue;
    }
    values.push(value);
  }

  return values;
}
