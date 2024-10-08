const scriptProp = PropertiesService.getScriptProperties()

function onOpen() {
  SpreadsheetApp.getUi()
    .createMenu("QazTrade")
    .addItem("Меню", "sideMenu")
    .addItem("Добавить компанию", "formAddCompany")
    .addItem("Добавить дистрибьюторский договор", "formAddDistributedContract")
    .addToUi();

  scriptProp.setProperty('key', SpreadsheetApp.getActiveSpreadsheet().getId())
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
    const htmlTemplate = HtmlService.createTemplateFromFile("form_edit");
    htmlTemplate.jsonBody = JSON.stringify(getSelectedRecord())

    SpreadsheetApp.getUi().showModalDialog(htmlTemplate.evaluate(), "Изменить запись");
}

function formAddCompany() {
    const html = HtmlService.createHtmlOutputFromFile("form_add_company");
    SpreadsheetApp.getUi().showModalDialog(html, "Добавить компанию");
}

function formAddDistributedContract() {
    const htmlTemplate = HtmlService.createTemplateFromFile("form_add_distributed_contract");
    htmlTemplate.parentValues = JSON.stringify(getParentValues("Дистрибьюторский договор"))

    SpreadsheetApp.getUi().showModalDialog(htmlTemplate.evaluate(), "Добавить дистрибьюторский договор");
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
        result += column.charAt(i).charCodeAt()  - 'A'.charCodeAt() + 1;
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
