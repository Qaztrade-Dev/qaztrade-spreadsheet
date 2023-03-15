class HeaderCell {
  constructor(Key, Range, Values, GroupKey = '') {
    this.Key = Key;
    this.Range = Range;
    this.Values = Values;
    this.GroupKey = GroupKey;
  }

  isLeaf() {
    return Object.keys(this.Values).length === 0;
  }
}

class Cell {
  constructor(value, rowNum, columnNum, headerCell) {
    this.value = value
    this.rowNum = rowNum
    this.columnNum = columnNum
    this.headerCell = headerCell
  }
}

const parents = {
  "№": {
      "parent": "root",
      "parentKey": "root"
  },
  "Дистрибьюторский договор": {
      "parent": "№",
      "parentKey": "№",
  },
  "контракт на поставку": {
      "parent": "Дистрибьюторский договор",
      "parentKey": "Дистрибьюторский договор.№",
  },
}

const children = {
  "root": [
    {
        "name": "№",
        "key": "№",
    }
  ],
  "№": [
      {
          "name": "Дистрибьюторский договор",
          "key": "Дистрибьюторский договор.№",
      }
  ],
  "Дистрибьюторский договор": [
      {
          "name": "контракт на поставку",
          "key": "контракт на поставку.№",
      }
  ],
  "контракт на поставку": [],
}

const headerCells = getHeaderCells();

function myFunc() {
  payload = {
    "parentID": null,
    "childKey": "№",
    "value": {
      "№": "1",
      "Производитель/дочерняя компания/дистрибьютор/СПК": "Doodocs",
      "подтверждающий документ": {
        "производитель": "Doodocs",
        "наименование": "Дудокс",
        "№": "3",
        "наименование товара": "Подписи",
        "ТН ВЭД (6 знаков)": "120934",
        "дата": "12.09.2019",
        "срок": "123",
        "подтверждение на сайте уполномоченного органа": "http://google.com",
      },
    }
  }

  rowNum = getRowNum(payload.parentID, payload.childKey)
  Logger.log("rowNum", rowNum)
  // insertRecord(payload.value, rowNum)
}

function getRowNum(parentID, childName) {
  // 1. get parent row
  // 2. get last child of the parent, e.g. neighbor
  // 3. get last row of the farthest descendent
  Logger.log("getRowNum: %s %s", parentID, childName)

  const parentKey = parents[childName].parentKey
  const parentBounds = getLevelBounds(parentKey, parentID)
  Logger.log("parentKey: %s", parentKey)
  Logger.log("parentBounds: %s", parentBounds)

  const upperBound = parentBounds[0]
  const lowerBound = parentBounds[1]
  
  const child = getChild(parentKey, childName)
  if (child == null) {
    return null
  }
  Logger.log("child: %s", child)
  const childHeaderCell = getHeaderCell(child.key)
  Logger.log("childHeaderCell: %s", childHeaderCell)
  const lastChildCell = getLastChildCell(parentBounds, childHeaderCell)
  Logger.log("lastChildCell: %s", lastChildCell)
  if (lastChildCell == null) {
    Logger.log("go here?")
    return upperBound+1
  }

  Logger.log("children: %s", children[lastChildCell.headerCell.Key][0])
  rowNum = getRowNum(lastChildCell.value, children[lastChildCell.headerCell.GroupKey][0].name)
  Logger.log("rowNum: %s", rowNum)
  if (rowNum == null) {
    return lastChildCell.rowNum
  }

  return rowNum

  // let sheet = SpreadsheetApp.getActiveSpreadsheet().getSheets()[0];
  // sheet.insertRowAfter(lastChildCell.RowNum)

  return lastChildCell.rowNum+1
}

function getLastChildCell(parentBounds, childHeaderCell) {
  const upperBound = parentBounds[0]+1
  const lowerBound = parentBounds[1]+1
  const columnNum = childHeaderCell.Range[0]

  let sheet = SpreadsheetApp.getActiveSpreadsheet().getSheets()[0];
  let rows = sheet.getRange(upperBound, columnNum, (lowerBound-upperBound)+1, 1).getValues()
  
  let lastIdx = 0
  let lastValue = ''

  for (i = 0; i < rows.length; i++) {
    const row = rows[i]
    const value = row[0]
    if (i == 0 && value === '') {
      return null
    }
    if (value !== '') {
      lastIdx = i
      lastValue = value
    }
  }

  return new Cell(lastValue.toString(), upperBound+lastIdx, columnNum, childHeaderCell)
}

function getChild(parentKey, childName) {
  if (!(parentKey in children)) {
    return null
  }

  for (i = 0; i < children[parentKey].length; i++) {
    const child = children[parentKey][i]
    if (child.name === childName) {
      return child
    }
  }
  return null
}

function getLevelBounds(parentKey, parentID) {
  // Logger.log("getLevelBounds: %s %s", parentKey, parentID)
  if (parentKey === "root") {
    // Logger.log("root getLevelBounds")
    var sheet = SpreadsheetApp.getActiveSpreadsheet().getSheets()[0];
    var range = sheet.getRange("A6:A");
    return [5, 5+range.getValues().length-1]
  }

  let cell = getHeaderCell(parentKey)
  // Logger.log("cell: %s", cell)

  let columnA1 = encodeAlphabet(cell.Range[0])
  // Logger.log("columnA1: %s", columnA1)
  
  var sheet = SpreadsheetApp.getActiveSpreadsheet().getSheets()[0];
  var range = sheet.getRange(columnA1+":"+columnA1);

  const rows = range.getValues()
  let l = 0;
  for (i = 0; i < rows.length; i++) {
    const row = rows[i]
    const rowValue = row[0]

    let value = rowValue
    if (typeof value === "number") {
      value = value.toString()
    }

    if (value === parentID) {
      l = i;
      break;
    }
  }

  let r = l;
  for (i = l+1; i < rows.length; i++) {
    const rowValue = rows[i][0]
    if (rowValue !== '') {
      break;
    }
    r = i
  }

  return [l, r];
}

function getHeaderCell(cellKey) {
  const keys = cellKey.split('.');
  let cell;

  for (let i = 0; i < keys.length; i++) {
    const key = keys[i]
    cell = (cell instanceof HeaderCell) ? cell.Values[key] : headerCells[key]
  }

  return cell;
}

function insertRecord(payload, rowNum) {
  fillSheet(payload, headerCells, rowNum);
}

function fillSheet(payload, headerCells, rowNum = 0) {
  var sheet = SpreadsheetApp.getActiveSpreadsheet().getSheets()[0];
  if (rowNum === 0) {
    rowNum = sheet.getLastRow() + 1;
  }

  Object.entries(payload).forEach(([k, v]) => {
    const cell = headerCells[k];
    if (cell.isLeaf()) {
      sheet.getRange(rowNum, cell.Range[0], 1, 1).setValues([[payload[k]]]);
    } else {
      fillSheet(payload[k], cell.Values, rowNum);
    }
  });
}

function getHeaderCells() {
  var sheet = SpreadsheetApp.getActiveSpreadsheet();
  var headerValues = sheet.getRangeByName("Header").getValues();

  topLevel = headerValues[0];
  lowLevel = headerValues[1];

  cellMap = {};
  for (i = 0; i < topLevel.length; i++) {
    if (topLevel[i] === "") {
      continue;
    }

    values = {};
    r = i;

    for (j = i; j < lowLevel.length; j++) {
      if (lowLevel[j] == "") {
        break;
      }
      if (!(topLevel[j] == "" || i == j)) {
        break;
      }
      values[lowLevel[j]] = new HeaderCell(lowLevel[j], [j + 1, j + 1], "", topLevel[i]);
      r = j;
    }

    cellMap[topLevel[i]] = new HeaderCell(topLevel[i], [i + 1, r + 1], values, topLevel[i]);
  }

  return cellMap;
}

function encodeAlphabet(num) {
  let result = "";
  while (num > 0) {
    let remainder = (num - 1) % 26;
    result = String.fromCharCode(65 + remainder) + result;
    num = Math.floor((num - 1) / 26);
  }
  return result;
}

function getColumnNum(column /* A,B,C,...,AA*/) {
  column = column.replace(/\d+/g, '');
  
  var result = 0;

  for (var i = 0; i < column.length; i++) {
      result *= 26;
      result += column.charAt(i).charCodeAt()  - 'A'.charCodeAt() + 1;
  }

  return result;
}


// input = {
    // "№": "1",
    // "Производитель/дочерняя компания/дистрибьютор/СПК": "Doodocs",
    // "подтверждающий документ": {
    //   "производитель": "Doodocs",
    //   "наименование": "Дудокс",
    //   "№": "3",
    //   "наименование товара": "Подписи",
    //   "ТН ВЭД (6 знаков)": "120934",
    //   "дата": "12.09.2019",
    //   "срок": "123",
    //   "подтверждение на сайте уполномоченного органа": "http://google.com",
    // },
  // };
  // payload = {
  //   "parentID": "1",
  //   "childKey": "Дистрибьюторский договор",
  //   "value": {
  //     "Дистрибьюторский договор": {
  //       "№": "1",
  //       "дата": "12.07.1997",
  //       "условия": "Салам всем!"
  //     }
  //   }
  // }