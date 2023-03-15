class Cell {
  constructor(Key, Range, Values) {
    this.Key = Key;
    this.Range = Range;
    this.Values = Values;
  }

  isLeaf() {
    return Object.keys(this.Values).length === 0;
  }
}

const parent = {
  "№": {
      "parent": null,
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

const child = {
  "№": [
      {
          "child": "Дистрибьюторский договор",
          "childKey": "Дистрибьюторский договор.№",
      }
  ],
  "Дистрибьюторский договор": [
      {
          "child": "контракт на поставку",
          "childKey": "контракт на поставку.№",
      }
  ],
  "контракт на поставку": [],
}

const headerCells = getHeaderCells();

function myFunc() {
  payload = {
    "parentID": "1",
    "childKey": "Дистрибьюторский договор",
    "value": {
      "Дистрибьюторский договор": {
        "№": "1",
        "дата": "12.07.1997",
        "условия": "Салам всем!"
      }
    }
  }

  rowNum = getRowNum(payload.parentID, payload.childKey)
  

  // insertRecord(input)
}

function getRowNum(parentID, childKey) {
  // 1. get parent row
  // 2. get last child of the parent, e.g. neighbor
  // 3. get last row of the farthest descendent

  parentKey = parent[childKey].parentKey
  parentBounds = getLevelBounds(parentKey, parentID)
  Logger.log(parentBounds)
}

function getLevelBounds(parentKey, parentID) {
  const cell = getCell(parentKey)
  const columnA1 = encodeAlphabet(cell.Range[0])
  
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

function getCell(cellKey) {
  const keys = cellKey.split('.');
  let cell;

  for (let i = 0; i < keys.length; i++) {
    cell = (typeof cell === "Cell") ? headerCells[keys[i]].Values : headerCells[keys[i]]
  }

  return cell;
}

function insertRecord(payload) {
  fillSheet(payload, headerCells);
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
      values[lowLevel[j]] = new Cell(lowLevel[j], [j + 1, j + 1], "");
      r = j;
    }

    cellMap[topLevel[i]] = new Cell(topLevel[i], [i + 1, r + 1], values);
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


// input = {
  //   "№": "1",
  //   "Производитель/дочерняя компания/дистрибьютор/СПК": "Doodocs",
  //   "подтверждающий документ": {
  //     "производитель": "Doodocs",
  //     "наименование": "Дудокс",
  //     "№": "3",
  //     "наименование товара": "Подписи",
  //     "ТН ВЭД (6 знаков)": "120934",
  //     "дата": "12.09.2019",
  //     "срок": "123",
  //     "подтверждение на сайте уполномоченного органа": "http://google.com",
  //   },
  // };