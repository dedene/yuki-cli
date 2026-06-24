package api

import (
	"context"
	"strings"
	"testing"
)

func TestUpdatedAndDeletedTransactionsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "ChangeDigest", "GetUpdatedAndDeletedTransactions", updatedAndDeletedTransactionsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:startDate>2025-07-23T00:00:00.00Z</they:startDate>",
			"<they:endDate>2025-08-23T13:00:00.00Z</they:endDate>",
			"<they:numberOfRecords>100</they:numberOfRecords>",
			"<they:startRecord>0</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	transactions, err := client.UpdatedAndDeletedTransactions(context.Background(), "session-1", UpdatedAndDeletedTransactionsOptions{
		AdministrationID: "admin-1",
		StartDate:        "2025-07-23T00:00:00.00Z",
		EndDate:          "2025-08-23T13:00:00.00Z",
		NumberOfRecords:  100,
		StartRecord:      0,
	})
	if err != nil {
		t.Fatalf("UpdatedAndDeletedTransactions: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("len(transactions) = %d, want 1", len(transactions))
	}
	tx := transactions[0]
	if tx.ID != "c6e77963-1eb6-4aea-bb65-77902c1b201d" ||
		tx.TransactionDate != "2023-02-06T00:00:00" ||
		tx.TransactionAmount != "-63.00" ||
		tx.Currency != "EUR" ||
		tx.ContactID != "64f5dfa0-d756-4fd2-8f1f-2b795439249a" ||
		tx.FullName != "Biovita Bvba" ||
		tx.GLAccountCode != "494100" ||
		tx.Updated != "2023-07-25T13:37:58.3" ||
		tx.Deleted != "true" {
		t.Fatalf("transaction = %#v", tx)
	}
}

func TestChangeDigestTransactionDetailParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "ChangeDigest", "GetTransactionDetail", changeDigestTransactionDetailResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:transactionID>tx-1</they:transactionID>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	tx, err := client.ChangeDigestTransactionDetail(context.Background(), "session-1", "admin-1", "tx-1")
	if err != nil {
		t.Fatalf("ChangeDigestTransactionDetail: %v", err)
	}
	if tx.ID != "1b477900-50c0-xxx-xxx-1de9ffe929c3" ||
		tx.HID != "919" ||
		tx.TransactionDate != "2021-01-01T00:00:00" ||
		tx.TransactionAmount != "0.00" ||
		tx.Currency != "EUR" ||
		tx.TaxCodeDescription != "TRMNon deductable" ||
		tx.FullName != "Topolino bvba" ||
		tx.ContactHID != "25" ||
		tx.ContactZipCode != "3000" ||
		tx.GLAccountCode != "494190" ||
		tx.DocumentID != "8f2fb08a-2c46-40d2-ab15-8758a869f32a" ||
		tx.ContactID != "5b0e9dc2-ec60-4614-afb6-d876a0f20370" ||
		tx.PeriodID != "202101" {
		t.Fatalf("transaction = %#v", tx)
	}
}

const updatedAndDeletedTransactionsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUpdatedAndDeletedTransactionsResponse xmlns="http://www.theyukicompany.com/">
      <GetUpdatedAndDeletedTransactionsResult>
        <UpdatedTransaction>
          <id>c6e77963-1eb6-4aea-bb65-77902c1b201d</id>
          <transactionDate>2023-02-06T00:00:00</transactionDate>
          <description>Factuur voor Biovita Bvba</description>
          <transactionAmount>-63.00</transactionAmount>
          <transactionAmountForeignCurrency>-63.00</transactionAmountForeignCurrency>
          <currencyRate>1.000000</currencyRate>
          <currency>EUR</currency>
          <contactID>64f5dfa0-d756-4fd2-8f1f-2b795439249a</contactID>
          <fullName>Biovita Bvba</fullName>
          <glAccountCode>494100</glAccountCode>
          <documentID />
          <created>2023-06-06T16:39:38.147</created>
          <updated>2023-07-25T13:37:58.3</updated>
          <deleted>true</deleted>
        </UpdatedTransaction>
      </GetUpdatedAndDeletedTransactionsResult>
    </GetUpdatedAndDeletedTransactionsResponse>
  </soap:Body>
</soap:Envelope>`

const changeDigestTransactionDetailResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetTransactionDetailResponse xmlns="http://www.theyukicompany.com/">
      <GetTransactionDetailResult>
        <id>1b477900-50c0-xxx-xxx-1de9ffe929c3</id>
        <hID>919</hID>
        <transactionDate>2021-01-01T00:00:00</transactionDate>
        <description>Factuur van Topolino bvba, Installaties, machines &amp; uitrusting</description>
        <transactionAmount>0.00</transactionAmount>
        <transactionAmountForeignCurrency>0.00</transactionAmountForeignCurrency>
        <currencyRate>1.000000</currencyRate>
        <currency>EUR</currency>
        <taxCodeDescription>TRMNon deductable</taxCodeDescription>
        <fullName>Topolino bvba</fullName>
        <CoCNumber>xxxxxxxx</CoCNumber>
        <VATNumber>BExxxxxxx</VATNumber>
        <contactHID>25</contactHID>
        <contactCountry>BE</contactCountry>
        <contactZipCode>3000</contactZipCode>
        <glAccountCode>494190</glAccountCode>
        <documentID>8f2fb08a-2c46-40d2-ab15-8758a869f32a</documentID>
        <documentReference>test</documentReference>
        <documentType>TRMPurchase invoice</documentType>
        <documentFolder>TRMPurchase_1</documentFolder>
        <documentFolderTab>TRMInvoices</documentFolderTab>
        <contactID>5b0e9dc2-ec60-4614-afb6-d876a0f20370</contactID>
        <periodId>202101</periodId>
        <company>la-cerva2</company>
        <mutationUser>yuki</mutationUser>
      </GetTransactionDetailResult>
    </GetTransactionDetailResponse>
  </soap:Body>
</soap:Envelope>`
