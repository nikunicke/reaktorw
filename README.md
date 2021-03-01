# reaktor warehouse v2

Here is my implementation for the *Reaktor Junior 2021 Summer* pre-assignment. Live version running on [**Heroku**](https://limitless-inlet-17222.herokuapp.com/). If the dyno has been a sleep, it will take about 15 seconds until data is available.

> *Your client is a clothing brand that is looking for a simple web app to use in their warehouses. To do their work efficiently, the warehouse workers need a fast and simple listing page per product category, where they can check simple product and availability information from a single UI.*

Read more about the assignment [**here**](https://www.reaktor.com/junior-dev-assignment/) *(27.2.2021)*


> :warning: **DISCLAIMER: This is not a simple solution to the assignment**. As this is my second iteration of the same project, I wanted to improve on what I had previously implemented and learn more about a few concepts that I recently became familiar with. Those would be *SOLID development principles, data-processing pipelines, service runners and asynchronous programming*. This being an assignment about abstracting data to an interpretable format for warehouse workers, I think those concepts fit really well. *(27.2.2021)*

---

*built with go1.15.8*

Based on the requirements of the assignment, this application should provide the following services:
*   A periodically running warehouse updater for keeping products and their availability status up to date by retrieving data from the provided API ([badapi](http://bad-api-assignment.reaktor.com/)), processing it and eventually storing it in the data warehouse. All requests to the API is executed in an asynchronous manner and for each manufacturer, multiple requests are sent to keep update times consistent.
*   A frontend for the end users to view products and their respective availability status.

The services are integreted into one application using a service runner, where each service is executed independently. The service runner keeps track of each service and exits gracefully if an error were to occur. 

The application is supported by the folloing packages:
* ### **Warehouse**
    *   Defines an inventory interface and can be implemented to support any DB management system. This project includes an implementation for a thread-safe in-memory store that allows only one read-write operation at a time but allows as many read-only transactions as you want at a time. Locks being quite slow, might cause bottlenecks during a warehouse update, as products and availability data is inserted in bulk. 
* ### **Pipeline**
    *   An asynchronous data-processing pipeline including two payload dispatch strategy types:
        *   FIFO: *First In, First Out*. The processor takes a payload, processes it and sends it to the next stage in an orderly manner.
        *   Fixed Worker Pool: If *out-of-order* processing is not an issue, a fixed worker pool can be used to process multiple payloads in parallel.
* ### **Badapi**
    *   A simple client for the badapi resource that includes iterator interfaces for adaptability with the pipeline package. The requests do not include Go contexts, meaning that the executed request is blocking until a response has been recieved or the client timeout has been exceeded.

