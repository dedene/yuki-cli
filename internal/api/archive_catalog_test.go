package api

import (
	"context"
	"strings"
	"testing"
)

func TestDocumentFoldersParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentFolders", documentFoldersResponse, nil)

	folders, err := client.DocumentFolders(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("DocumentFolders: %v", err)
	}
	if len(folders) != 3 {
		t.Fatalf("len(folders) = %d, want 3", len(folders))
	}
	if folders[0].ID != "7" ||
		folders[0].Description != "To be handled by Yuki" ||
		folders[0].Icon != "DocumentFolder_yellow_label.png" ||
		!folders[0].ProcessedByYuki ||
		folders[2].ID != "2" ||
		folders[2].Description != "Sales" {
		t.Fatalf("folders = %#v", folders)
	}
}

func TestDocumentFolderTabsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentFolderTabs", documentFolderTabsResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:folderID>3</they:folderID>") {
			t.Fatalf("request body missing folder ID:\n%s", body)
		}
	})

	tabs, err := client.DocumentFolderTabs(context.Background(), "session-1", "3")
	if err != nil {
		t.Fatalf("DocumentFolderTabs: %v", err)
	}
	if len(tabs) != 5 {
		t.Fatalf("len(tabs) = %d, want 5", len(tabs))
	}
	if tabs[0].ID != "301" ||
		tabs[0].Description != "Files" ||
		!tabs[0].ProcessedByYuki ||
		tabs[2].ID != "303" ||
		tabs[2].Description != "Credit cards" {
		t.Fatalf("tabs = %#v", tabs)
	}
}

func TestPaymentMethodsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "PaymentMethods", archivePaymentMethodsResponse, nil)

	methods, err := client.PaymentMethods(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("PaymentMethods: %v", err)
	}
	if len(methods) != 2 ||
		methods[0].ID != "4" ||
		methods[0].Description != "Zakelijke Bancontact" ||
		methods[1].ID != "5" ||
		methods[1].Description != "Zakelijke Credit card" {
		t.Fatalf("methods = %#v", methods)
	}
}

func TestCurrenciesParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "Currencies", currenciesResponse, nil)

	currencies, err := client.Currencies(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("Currencies: %v", err)
	}
	if len(currencies) != 4 {
		t.Fatalf("len(currencies) = %d, want 4", len(currencies))
	}
	if currencies[0].ID != "EUR" ||
		!currencies[0].Default ||
		currencies[0].Description != "Euro (EUR)" ||
		currencies[3].ID != "USD" ||
		currencies[3].Default {
		t.Fatalf("currencies = %#v", currencies)
	}
}

func TestCostCategoriesParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "CostCategories", costCategoriesResponse, nil)

	categories, err := client.CostCategories(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("CostCategories: %v", err)
	}
	if len(categories) != 2 {
		t.Fatalf("len(categories) = %d, want 2", len(categories))
	}
	if categories[0].ID != "40300" ||
		categories[0].Description != "Training costs" ||
		categories[1].ID != "40600" ||
		categories[1].Description != "Canteen supplies" {
		t.Fatalf("categories = %#v", categories)
	}
}

func TestMenuParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "Menu", menuResponse, nil)

	entries, err := client.Menu(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("Menu: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("len(entries) = %d, want 3", len(entries))
	}
	if entries[0].ID != "1" ||
		entries[0].Text != "Vragen" ||
		entries[0].Icon != "yw3-question" ||
		entries[1].Alert != "314" ||
		entries[2].Link != "IPPLReport.aspx" {
		t.Fatalf("entries = %#v", entries)
	}
}

const archivePaymentMethodsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <PaymentMethodsResponse xmlns="http://www.theyukicompany.com/">
      <PaymentMethodsResult>
        <PaymentMethods xmlns="">
          <PaymentMethod ID="4">
            <Description>Zakelijke Bancontact</Description>
          </PaymentMethod>
          <PaymentMethod ID="5">
            <Description>Zakelijke Credit card</Description>
          </PaymentMethod>
        </PaymentMethods>
      </PaymentMethodsResult>
    </PaymentMethodsResponse>
  </soap:Body>
</soap:Envelope>`

const documentFoldersResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentFoldersResponse xmlns="http://www.theyukicompany.com/">
      <DocumentFoldersResult>
        <DocumentFolders xmlns="">
          <DocumentFolder ID="7">
            <Description>To be handled by Yuki</Description>
            <Icon>DocumentFolder_yellow_label.png</Icon>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolder>
          <DocumentFolder ID="1">
            <Description>Purchase</Description>
            <Icon>DocumentFolder_red_label.png</Icon>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolder>
          <DocumentFolder ID="2">
            <Description>Sales</Description>
            <Icon>DocumentFolder_red_label.png</Icon>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolder>
        </DocumentFolders>
      </DocumentFoldersResult>
    </DocumentFoldersResponse>
  </soap:Body>
</soap:Envelope>`

const documentFolderTabsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentFolderTabsResponse xmlns="http://www.theyukicompany.com/">
      <DocumentFolderTabsResult>
        <DocumentFolderTabs xmlns="">
          <DocumentFolderTab ID="301">
            <Description>Files</Description>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolderTab>
          <DocumentFolderTab ID="302">
            <Description>Statement view</Description>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolderTab>
          <DocumentFolderTab ID="303">
            <Description>Credit cards</Description>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolderTab>
          <DocumentFolderTab ID="304">
            <Description>Petty cash</Description>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolderTab>
          <DocumentFolderTab ID="305">
            <Description>Other</Description>
            <ProcessedByYuki>True</ProcessedByYuki>
          </DocumentFolderTab>
        </DocumentFolderTabs>
      </DocumentFolderTabsResult>
    </DocumentFolderTabsResponse>
  </soap:Body>
</soap:Envelope>`

const currenciesResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CurrenciesResponse xmlns="http://www.theyukicompany.com/">
      <CurrenciesResult>
        <Currencies xmlns="">
          <Currency ID="EUR" Default="True">
            <Description>Euro (EUR)</Description>
          </Currency>
          <Currency ID="GBP">
            <Description>British pound</Description>
          </Currency>
          <Currency ID="ISK">
            <Description>Icelandic króna</Description>
          </Currency>
          <Currency ID="USD">
            <Description>US dollar</Description>
          </Currency>
        </Currencies>
      </CurrenciesResult>
    </CurrenciesResponse>
  </soap:Body>
</soap:Envelope>`

const costCategoriesResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CostCategoriesResponse xmlns="http://www.theyukicompany.com/">
      <CostCategoriesResult>
        <CostCategories xmlns="">
          <CostCategory ID="40300">
            <Description>Training costs</Description>
          </CostCategory>
          <CostCategory ID="40600">
            <Description>Canteen supplies</Description>
          </CostCategory>
        </CostCategories>
      </CostCategoriesResult>
    </CostCategoriesResponse>
  </soap:Body>
</soap:Envelope>`

const menuResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <MenuResponse xmlns="http://www.theyukicompany.com/">
      <MenuResult>
        <Menu xmlns="">
          <MenuEntry ID="1">
            <Text>Vragen</Text>
            <Icon>yw3-question</Icon>
            <Link>IPQuestions.aspx</Link>
            <Alert>0</Alert>
          </MenuEntry>
          <MenuEntry ID="2">
            <Text>Aandacht</Text>
            <Icon>yw3-alert</Icon>
            <Link>IPAlert.aspx</Link>
            <Alert>314</Alert>
          </MenuEntry>
          <MenuEntry ID="3">
            <Text>Resultaten</Text>
            <Icon>yw3-chart</Icon>
            <Link>IPPLReport.aspx</Link>
            <Alert>0</Alert>
          </MenuEntry>
        </Menu>
      </MenuResult>
    </MenuResponse>
  </soap:Body>
</soap:Envelope>`
