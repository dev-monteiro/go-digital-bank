# Go Account Fee

This is a personal project to practice what I've learned about Golang so far (and learn a lot more). The goal is to implement some enterprise-like features in a fictional digital bank context in a way that I will be able to:
- have a taste of what is like working with Go on a real-life project;
- build common components of modern cloud systems, like microservices and serverless functions;
- explore the benefits and limitations of working with Golang and its tooling.

## Context

Over the last decade, it has been seen a boom of digital banks popping up in Brazil (for the record, I currently work in one of them). Like most Brazilians, my first experience with these banks was being a customer of Nubank, the most popular of the sort in Latin America, which has arguably the best app UI/UX, being miles ahead of its competitors even today. One of the things that I liked about this app is that the credit card invoice is updated almost in real-time, regardless if you are using it in some local establishment or paying some amount of the invoice with the digital account balance. In contrast, most of the other banks (especially the ones that were there before the digital era) take from 1 to 2 days to make this kind of update; some of them don't even show which transactions are pending, in a way that they just pop up on the invoice next day. 

I think that this behavior results from all the complicated steps to finish a credit card transaction, in a way that the invoice is only updated when everything is finished (maybe there is even a regulation that says that should be like that). I know that the born-digital ones usually are built over some off-the-shelf core banking system, which, in turn, has these daily batch processes to settle transactions, with the invoices being updated afterward. So, I don't know what Nubank does, but I think that if one of those digital banks wanted to reproduce this behavior, it would probably assume, for the current day, an approved transaction as settled, in the sense that they can be updated it right away in the invoice view for the user. This way, the transaction would be settled later on the batch processing and included in the real invoice, which, in turn, would replace the temporary view for the user and, on the happy path, this change of views won't be noticed and the user will feel like the invoice is updated immediately. Of course, it must have some (or a lot) of pitfalls in this approach, but I think that it can be an interesting challenge for a practice project with the goals that I listed before, so that's exactly I will be trying to do here.

In resume, the challenge is to develop a credit card invoice system that seems to provide immediately updated info to the user in the scenario that this development is done in a new-born digital bank which is built over some off-the-shelf core banking system.

## Business

Lorem ipsum.


## Front-End Contract

Below, we can see screens from Nubank to displaying invoices info, in which we'll just worry about the info on the highlighted areas. The left and center ones refer to the current invoice screen in its two states: closed and open. When tapping on the invoices summary button, the user is forwarded to the right screen, which starts showing the current invoice entries and can be swiped to select the other ones.

{screens}

Our system needs to provides an API that allows the building of these screens (at least for the highlighted areas). This way, for the current invoice info, we can define the following endpoint:

Endpoint
```
GET {host}/invoices/current?customerId={customerId}
```

Response Body
```json
{
    "id": "abc-123-def",
    "statusLabel": "Closed|Current",
    "amount": 1234.56,
    "closingDate": "APR 13"
}
```

As for the invoice summary, we can define the following one. It worth notices that this screen sample doesn't have an example of a transaction with installments, in which the entry amount would refer to that month installment and additional info would append the entry description in this way: "Some purchase with installments  3/6".

Endpoint
```
GET {host}/invoices/summary/{invoiceId}
```

Response Body
```json
{
    "id": "abc-123-def",
    "refMonthLabel": "NOV",
    "amount": 1675.55,
    "dueDate": "NOV 8",
    "closingDate": "NOV 1",
    "entries": [
        {
            "date": "OCT 24",
            "description": "Taxi SP",
            "amount": 56.54
        },
        {
            "date": "OCT 27",
            "description": "Amazon 3/6",
            "amount": 111.11
        }
    ]
}
```


## Core Banking Contract

As mentioned, the development should be made on top of some off-the-shelf core banking system, so let's define the API and events contracts for the latter. we'll take as reference the documentation provided by [Dock](https://lighthouse.dock.tech), a common vendor of this kind of system in Brazil. Also, for simplicity, let's not worry about some common API aspects like authentication and pagination for now.

First, we need a way to query the invoices settled by the core banking system for a given customer. For that, we'll use the API defined below based on [this one](https://lighthouse.dock.tech/docs/pier-pro-api-reference/1b0ec155a9966-list-of-invoices):

Endpoint
```
GET {host}/invoices?customerId=:customerId
```

Response
```json
{
    "customerId": 123,
    "invoiceId": 1234,
    "processingSituation": "CLOSED|OPEN|FUTURE",
    "amount": 1234.56,
    "closingDate": "2023-08-05",
    "dueDate": "2023-08-15",
    "actualDueDate": "2023-08-15",
    "isPaymentDone": true
}
```

Then, we need a way to query a specific invoice and get all its entries. For that, we'll use the API defined below based on [this one](https://lighthouse.dock.tech/docs/pier-pro-api-reference/974ef197cb242-retrieve-the-invoice-of-a-client):

Endpoint
```
GET {host}/invoices/:invoiceId?customerId=:customerId
```

Response
```json
{
    "customerId": 123,
    "invoiceId": 1234,
    "processingSituation": "CLOSED|OPEN|FUTURE",
    "amount": 1234.56,
    "closingDate": "2023-08-05",
    "dueDate": "2023-08-15",
    "actualDueDate": "2023-08-15",
    "isPaymentDone": true,
    "entries": [
        {
            "transactionId": 12345,
            "amount": 22.01,
            "installmentNumber": 2,
            "numInstallments": 3,
            "transactionDateTime": "2023-05-31T09:54:30",
            "merchantDescription": "DFV Digital"
        }
    ]
}
```

Also, we need to receive transactional events happening throughout the day. The documentation describes [some types of them](https://lighthouse.dock.tech/docs/pier-pro-api-reference/1729e0328e134-events#purchase-events) (seemed to be divided by the processing step), but they are very similar, so we're gonna use only one event model and distinct them by status:

Event
```json
{
    "transactionId": 12345,
    "customerId": 123,
    "transactionDateTime": "2023-05-31T09:54:30",
    "totalAmount": 66.03,
    "numInstallments": 3,
    "merchantDescription": "DFV Digital",
    "status": "PENDING|CLEARED|PROCESSED|CANCELED"
}
```

Finally, we need some way to know if the daily batches were processed so that we can trust the APIs to updated regarding the previous day. The documentation do have an event model for the batch processing but I couldn't find what is the link between the batches and the customers. For that, to add some logic complexity, I'll assume that every batch is meant to process a bunch of credit cards.