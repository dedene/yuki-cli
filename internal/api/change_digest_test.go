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
