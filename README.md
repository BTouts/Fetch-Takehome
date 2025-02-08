Fetch Take-Home Assignment Submission

This repository contains my implementation of the receipt processing web service for the Fetch take-home interview. The service follows the provided specifications and allows users to submit receipts and retrieve the associated points based on the defined scoring rules.

1. Clone the repository and cd into the project folder
2. To install dependencies, run: go mod tidy
3. To start the webservice, run: go run .
4. By default, the server will be available at http://localhost:8080
5. To process a receipt, use Postman, cURL, or any API testing tool on the following endpoint:
     POST /receipts/process
6. To retrieve the points for a receipt, use the below endpoint:
     GET /receipts/{id}/points
7. A small Postman collection is included in the testing folder that stores the most recently generated id from the POST request to make it a little easier to test with.
8. 
   
