# READ ME
## Building the container
From within your gcloud environment, run the following terminal command inside the project folder:

`gcloud builds submit --tag gcr.io/rockgymbuddy/report-service .`

Information on Cloud Build Pricing : https://cloud.google.com/cloud-build/Pricing

## Setup BigQuery in your project
For more information on setting up bigquery, follow the guide here : https://cloud.google.com/bigquery/docs/quickstarts/quickstart-command-line#step_2_create_a_new_dataset
### Add a new dataset
`bq mk reports`

### Add a new table within that dataset
`bq mk -t --description "Where we will collect reports" reports.reporttable cpu:STRING,memory:STRING,storage:STRING,deviceid:STRING`

You are now able to view the table from the dashboard and are configured for streaming inserts.

## Important environment variables to set in Cloud Run
Also be sure to set your cloud run environment to allow unauthenticated requests.
* `PROJ_ID` - The id of your project as described by Google
* `DATASET_ID` - The BigQuery dataset id set in previous step
* `TABLE_ID` - The table specified in BigQuery via the previous step

## Make post requests
Once the container is deployed, you can start making post requests to the container.  From here, you should be good to start sending the metrics you are interested in recording over to BigQuery.  After that, you can export the data to datastudio or query the data within the webs bigquery dashboard : https://console.cloud.google.com/bigquery
