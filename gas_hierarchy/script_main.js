function onOpen() {
  refreshMenu();
}

function refreshMenu() {
  name = "QazTrade";
  SpreadsheetApp.getActiveSpreadsheet().removeMenu(name);
  createMenu(name);
}

function createMenu(name) {
  menu = SpreadsheetApp.getUi();

  menu = menu
    .createMenu(name)
    .addItem("Заполнить заявление", "modalApplication")
    .addSubMenu(
      SpreadsheetApp.getUi()
        .createMenu("Доставка ЖД транспортом")
        .addItem("Добавить Лист", "trainDeliveryHandler")
        .addItem('Добавить "Компанию"', "trainDeliveryCompany")
        .addItem(
          'Добавить "Дистрибьюторский договор"',
          "trainDeliveryDistributedContractAdd"
        )
        .addItem(
          'Изменить "Дистрибьюторский договор"',
          "trainDeliveryDistributedContractEdit"
        )
    )
    .addSubMenu(
      SpreadsheetApp.getUi()
        .createMenu("Затраты на продвижение")
        .addItem("Добавить Лист", "advertisementExpensesHandler")
        .addItem('Добавить "Начало заявки"', "advertisementExpensesBeginning")
    );

  menu.addToUi();
}

const apiHost = "https://0161-2a09-bac5-47fb-b05-00-119-d.eu.ngrok.io";

function modalApplication() {
  const htmlTemplate = HtmlService.createTemplateFromFile("html_application");
  htmlTemplate.applicationObj = getEncodedApplicationValues();

  SpreadsheetApp.getUi().showModalDialog(htmlTemplate.evaluate(), "Заявление");
}

function trainDeliveryHandler() {
  sheetName = "Доставка ЖД транспортом";
  if (hasSheet(sheetName)) {
    return;
  }

  let options = {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify({ sheet_name: "Доставка ЖД транспортом" }),
    headers: {
      Authorization: "bearer " + getToken(),
    },
  };

  UrlFetchApp.fetch(`${apiHost}/sheets/`, options);
  refreshMenu();
}

function trainDeliveryCompany() {
  sheetName = "Доставка ЖД транспортом";
  if (!hasSheet(sheetName)) {
    return;
  }

  const htmlTemplate = HtmlService.createTemplateFromFile(
    "html_train_delivery_company"
  );
  htmlTemplate.token = getToken();
  htmlTemplate.sheet_name = sheetName;
  htmlTemplate.sheet_id = getSheetId(sheetName);

  SpreadsheetApp.getUi().showModalDialog(htmlTemplate.evaluate(), "Заявление");
}

function trainDeliveryDistributedContractAdd() {
  sheetName = "Доставка ЖД транспортом";
  if (!hasSheet(sheetName)) {
    return;
  }

  const htmlTemplate = HtmlService.createTemplateFromFile(
    "html_train_delivery_distributed_contract"
  );
  htmlTemplate.token = getToken();
  htmlTemplate.sheet_name = sheetName;
  htmlTemplate.sheet_id = getSheetId(sheetName);
  htmlTemplate.parentValues = JSON.stringify(
    getParentValues(sheetName, "Дистрибьюторский договор")
  );
  htmlTemplate.rowValues = JSON.stringify([]);
  htmlTemplate.rowNum = 0;

  SpreadsheetApp.getUi().showModalDialog(htmlTemplate.evaluate(), "Заявление");
}

function trainDeliveryDistributedContractEdit() {
  sheetName = "Доставка ЖД транспортом";
  if (!hasSheet(sheetName)) {
    return;
  }

  const htmlTemplate = HtmlService.createTemplateFromFile(
    "html_train_delivery_distributed_contract"
  );
  htmlTemplate.token = getToken();
  htmlTemplate.sheet_name = sheetName;
  htmlTemplate.sheet_id = getSheetId(sheetName);
  htmlTemplate.parentValues = JSON.stringify([]);
  htmlTemplate.rowValues = JSON.stringify(getRowValues());
  htmlTemplate.rowNum = getCurrentRowNum();

  SpreadsheetApp.getUi().showModalDialog(htmlTemplate.evaluate(), "Заявление");
}

function advertisementExpensesHandler() {
  sheetName = "Затраты на продвижение";
  if (hasSheet(sheetName)) {
    return;
  }

  let options = {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify({ sheet_name: sheetName }),
    headers: {
      Authorization: "bearer " + getToken(),
    },
  };

  UrlFetchApp.fetch(`${apiHost}/sheets/`, options);
  refreshMenu();
}

function advertisementExpensesBeginning() {
  sheetName = "Затраты на продвижение";
  if (!hasSheet(sheetName)) {
    return;
  }

  const htmlTemplate = HtmlService.createTemplateFromFile(
    "html_advertisement_expenses_beginning"
  );
  htmlTemplate.token = getToken();
  htmlTemplate.sheet_name = sheetName;
  htmlTemplate.sheet_id = getSheetId(sheetName);
  htmlTemplate.parentValues = JSON.stringify([]);
  htmlTemplate.rowValues = JSON.stringify([]);
  htmlTemplate.rowNum = 0;

  SpreadsheetApp.getUi().showModalDialog(htmlTemplate.evaluate(), sheetName);
}

function getSheetId(name) {
  const sheets = SpreadsheetApp.getActiveSpreadsheet().getSheets();
  for (i = 0; i < sheets.length; i++) {
    if (sheets[i].getName() == name) {
      return sheets[i].getSheetId();
    }
  }

  return null;
}

function hasSheet(name) {
  const sheets = SpreadsheetApp.getActiveSpreadsheet().getSheets();
  for (i = 0; i < sheets.length; i++) {
    if (sheets[i].getName() == name) {
      return true;
    }
  }

  return false;
}

function getEncodedApplicationValues() {
  applicationObj = getApplicationValues();
  applicationObj["token"] = getToken();

  const queryParams = Object.keys(applicationObj)
    .map((key) => {
      return `${key}=${encodeURIComponent(applicationObj[key])}`;
    })
    .join("&");

  return queryParams;
}

function getApplicationValues() {
  applicationSheet =
    SpreadsheetApp.getActiveSpreadsheet().getSheetByName("Заявление");
  namedRanges = applicationSheet.getNamedRanges();
  applicationObj = {};

  for (i = 0; i < namedRanges.length; i++) {
    key = namedRanges[i].getName();
    value = namedRanges[i].getRange().getValues()[0][0];
    applicationObj[key] = value;
  }

  return applicationObj;
}

function getToken() {
  key = "token";
  return getMetadataValue(key);
}

function getMetadataValue(key) {
  metadataFinder =
    SpreadsheetApp.getActiveSpreadsheet().createDeveloperMetadataFinder();
  metadata = metadataFinder.withKey(key).find();
  return metadata[metadata.length - 1].getValue();
}
