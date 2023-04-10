class HeaderCell {
  constructor(key, range, values, groupKey = "") {
    this.key = key;
    this.range = range;
    this.values = values;
    this.groupKey = groupKey;
  }
}

function getHeaderCells(sheetName) {
  const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
  const tmpRangeName = `${sheetName}!${sheetName.replaceAll(" ", "_")}_header`;
  const headerValues = spreadsheet.getRangeByName(tmpRangeName).getValues();

  const topLevel = headerValues[0];
  const lowLevel = headerValues[1];

  const cellMap = {};
  for (let i = 0; i < topLevel.length; i++) {
    if (topLevel[i] === "") {
      continue;
    }

    const values = {};
    let r = i;

    for (let j = i; j < lowLevel.length; j++) {
      if (lowLevel[j] == "") {
        break;
      }

      if (!(topLevel[j] == "" || i == j)) {
        break;
      }

      values[lowLevel[j]] = new HeaderCell(
        lowLevel[j],
        [j + 1, j + 1],
        {},
        topLevel[i]
      );
      r = j;
    }

    cellMap[topLevel[i]] = new HeaderCell(
      topLevel[i],
      [i + 1, r + 1],
      values,
      topLevel[i]
    );
  }

  return cellMap;
}

function getHeaderCell(sheetName, cellKey) {
  const headerCells = getHeaderCells(sheetName);
  const keys = cellKey.split("|");
  let cell = headerCells[keys[0]];

  for (let i = 1; i < keys.length; i++) {
    const key = keys[i];
    cell = cell.values[key];
  }

  return cell;
}

const trainDeliveryParents = {
  "№": {
    parent: "root",
    parentKey: "root",
  },
  "Дистрибьюторский договор": {
    parent: "№",
    parentKey: "№",
  },
  "контракт на поставку": {
    parent: "Дистрибьюторский договор",
    parentKey: "Дистрибьюторский договор|№",
  },
};

const advertisementExpensesParents = {
  "№": {
    parent: "root",
    parentKey: "root",
  },
};

const sheetParents = {
  "Доставка ЖД транспортом": trainDeliveryParents,
  "Затраты на продвижение": advertisementExpensesParents,
};

// function test() {
//   sheetName = 'Доставка ЖД транспортом'
//   childKey = 'Дистрибьюторский договор'

//   values = getParentValues(sheetName, childKey)
//   Logger.log(values)
// }

function getParentValues(sheetName, childKey) {
  const parentHeaderCell = getHeaderCell(
    sheetName,
    sheetParents[sheetName][childKey].parentKey
  );
  const values = getLevelValues(sheetName, parentHeaderCell);

  const filteredValues = values.reduce((nonEmptyStrings, str) => {
    if (typeof str === "string" && str.trim() !== "") {
      nonEmptyStrings.push(str);
    }
    return nonEmptyStrings;
  }, []);

  return filteredValues;
}

function getLevelValues(sheetName, headerCell) {
  const upperBound = 3;
  const columnNum = headerCell.range[0];

  const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(sheetName);
  const rows = sheet
    .getRange(upperBound, columnNum, sheet.getMaxRows(), 1)
    .getValues();

  const values = [];
  for (let i = 0; i < rows.length; i++) {
    const value = rows[i][0].toString();
    if (value === "") {
      continue;
    }
    values.push(value);
  }

  return values;
}

function getRowValues() {
  var spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
  var sheet = spreadsheet.getActiveSheet();
  // var sheet = spreadsheet.getSheets()[1];
  var activeCell = sheet.getActiveCell();

  var rowNum = activeCell.getRowIndex();
  // var rowNum = 5;
  var rowValues = sheet
    .getRange(rowNum, 1, 1, sheet.getLastColumn())
    .getValues()[0];

  const headerCellsMap = getHeaderCells(sheet.getName());
  const sortedHeaderCellsMap = sortHeaderCells(headerCellsMap);
  const rowValue = mapRowValues(sortedHeaderCellsMap, rowValues);

  return rowValue;
}

function getCurrentRowNum() {
  var spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
  var sheet = spreadsheet.getActiveSheet();
  var activeCell = sheet.getActiveCell();

  const sheetName = sheet.getName();
  const dataRangeName = `${sheetName}!${sheetName.replaceAll(" ", "_")}_data`;
  const dataRowNum = spreadsheet.getRangeByName(dataRangeName).getRowIndex();

  var rowNum = activeCell.getRowIndex() - dataRowNum;

  return rowNum;
}

function getRowNum(record) {
  record = record.replace(/\D/g, "");
  return parseInt(record);
}

function sortHeaderCells(cellMap) {
  const arr = [];
  for (var k in cellMap) {
    arr.push(cellMap[k]);
  }

  arr.sort((a, b) => (a.range[0] > b.range[0] ? 1 : -1));

  return arr;
}

function mapRowValues(headerCellsMap, rowValues) {
  let result = {};

  for (let headerCell of headerCellsMap) {
    let key = headerCell.key;
    let range = headerCell.range;

    if (headerCell.values && Object.keys(headerCell.values).length > 0) {
      result[key] = {};

      for (let subKey in headerCell.values) {
        let subRange = headerCell.values[subKey].range;
        result[key][subKey] = rowValues[subRange[0] - 1];
      }
    } else {
      result[key] = rowValues[range[0] - 1];
    }
  }

  return result;
}
