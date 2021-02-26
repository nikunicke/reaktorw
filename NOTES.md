# Notes - Reaktor Warehouse

* Assignment: https://www.reaktor.com/junior-dev-assignment/

## BAD-API
* Base URL:
    *   v1: https://bad-api-assignment.reaktor.com/
    *   v2: https://bad-api-assignment.reaktor.com/v2/


| Request Type | Path | Description | Format | Load(ms) |
|--------------|------|-------------|--------|------|
| GET | /products/{category} | Get all products in specified category | application/json | ~500ms |
| GET | /availability/{manufacturer} | Get availability data for specified manufacturer | application/json | ~15000ms |

In addition to high latency, the */availability/{manufacturer}* -endpoint might respond with partial data. In case of this error, the response header `X-Error-Modes-Active` is set. Inspection of the response body would still be recommended.

The data from the bad-api seems to be completely flushed and repopulated a few times per day. This could be interpreted in two ways: either the bad-api endpoints respond with data that was most recently updated , or, in the context of this being a pre-assignment, the data is just overwritten with randomly generated content. By tracking some products and finding out they never reappear (get updated) we can conclude that the latter is the most likely scenario.

## Approach
*   Asynchronous programming
    *   The bad-api problem can be approched by simply making asynchronous requests until the response contains the expected data. This adds complexity to our program as we need to deal with the following:
        *   Manage the lifecycle of each process
        *   Make use of concurrency tools, such as locks, to avoid data races

        Nevertheless, the asynchronous approach will result in consistent processing times compared to a synchronous approach. Even developing and testing will feel like less of a hazzle.
*   Simple data warehouse
    *   Simple thread-safe in-memory data management package that can categorize products and manufacturers.
    *   Might need to be flushed completely everytime we update. Certainly tough on the GC but probably the easiest solution. Optionally we iterate all keys/values and delete.
    *   Generating new UUIDs for the products as the bad-api is not quite reliable in any way...
*   gRPC
    *   This is mostly for learning purposes but could actually be quite a good choice for this assignment as it is very convenient for transfering a lot of data. And to top it off, server streaming would be cool for live data.
*   Frontend
    *   Frontend to display a table with catetgorized prodcts and their availability status. Optionally add some visualization for manufacturer overall availability status.
        *   Bootstrap and D3
        *   Google Analytics
        *   BI tool
