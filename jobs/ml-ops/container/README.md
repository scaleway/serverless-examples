# Deploy an inference API on Scaleway Serverless Containers

## Test the inference API using HTTP calls

You can perform the following HTTP calls:

```bash
curl -H "X-Auth-Token: $CONTAINER_TOKEN" -X POST "<scw_container_endpoint>" -H "Content-Type: application/json" -d '{"age": 44, "job": "blue-collar", "marital": "married", "education": "basic.4y", "default": "unknown", "housing": "yes", "loan": "no", "contact": "cellular", "month": "aug", "day_of_week": "thu", "duration": 210, "campaign": 1, "pdays": 999, "previous": "0", "poutcome": "nonexistent", "emp_var_rate": 1.4, "cons_price_idx": 93.444, "cons_conf_idx": -36.1, "euribor3m": 4.963, "nr_employed": 5228.1}'  
```

In order to:

* Load the model after training using the `/load_classifier`.
* Call inference endpoint using `/inference`.