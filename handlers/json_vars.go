package handlers

var ESTIMATE string = `{
  "Line": [
    {
      "Description": "Pest Control Services",
      "DetailType": "SalesItemLineDetail",
      "SalesItemLineDetail": {
        "TaxCodeRef": {
          "value": "NON"
        },
        "Qty": 1,
        "UnitPrice": 35,
        "ItemRef": {
          "name": "Pest Control",
          "value": "10"
        }
      },
      "LineNum": 1,
      "Amount": 35.0,
      "Id": "1"
    },
    {
      "DetailType": "SubTotalLineDetail",
      "Amount": 35.0,
      "SubTotalLineDetail": {}
    },
    {
      "DetailType": "DiscountLineDetail",
      "Amount": 3.5,
      "DiscountLineDetail": {
        "DiscountAccountRef": {
          "name": "Discounts given",
          "value": "86"
        },
        "PercentBased": true,
        "DiscountPercent": 10
      }
    }
  ],
  "CustomerRef": {
    "name": "Cool Cars",
    "value": "3"
  },
  "TxnTaxDetail": {
    "TotalTax": 0
  },
  "ApplyTaxAfterDiscount": false
}`

var INVOICE string = `{
  "Line": [
    {
      "DetailType": "SalesItemLineDetail",
      "Amount": 100.0,
      "SalesItemLineDetail": {
        "ItemRef": {
          "name": "Services",
          "value": "1"
        }
      }
    }
  ],
  "CustomerRef": {
    "value": "2"
  }
}`
