function onOpen() {
  SpreadsheetApp.getUi()
    .createMenu("QazTrade")
    .addItem("Меню", "sideMenu")
    .addItem("Добавить запись", "formAdd")
    .addItem("Изменить запись", "formEdit")
    .addToUi();
}

function sideMenu() {
    const html = HtmlService.createHtmlOutputFromFile("side_menu");
    html.setTitle("QazTrade меню");
    SpreadsheetApp.getUi().showSidebar(html);
}

function formAdd() {
    const html = HtmlService.createHtmlOutputFromFile("form_add");
    SpreadsheetApp.getUi().showModalDialog(html, "Добавить запись");
}

function formEdit() {
    const html = HtmlService.createHtmlOutputFromFile("form_edit");
    SpreadsheetApp.getUi().showModalDialog(html, "Изменить запись");
}

function getSelectedRecord() {
    var activeSheet = SpreadsheetApp.getActiveSheet();
    var selectedCell = activeSheet.getSelection().getCurrentCell().getA1Notation();
    var rowNum = getRowNum(selectedCell)

    var rowValues = activeSheet.getRange(rowNum, 1, 1, activeSheet.getLastColumn()).getValues()[0]
    var headers = activeSheet.getRange(1, 1, 1, activeSheet.getLastColumn()).getValues()[0]

    const jsonValues = convertToObj(headers, rowValues)
    jsonValues.rowNum = rowNum

    return jsonValues
}

function getColumnNum(column /* A,B,C,...,AA*/) {
    column = column.replace(/\d+/g, '');
    
    var result = 0;

    for (var i = 0; i < column.length; i++) {
        result *= 26;
        result += c  - 'A'.charCodeAt() + 1;
    }
 
    return result;
}

function getRowNum(record) {
    record = record.replace(/\D/g,'');
    return parseInt(record);
}

function convertToObj(a, b){
    if(a.length != b.length || a.length == 0 || b.length == 0){
        return null;
    }
    let obj = {};
     
    // Using the foreach method
    a.forEach((k, i) => {obj[k] = b[i]})
    return obj;
}
