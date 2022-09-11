# Visual-Programming-api

This api is responsible for storing the structure of the programs, in a [Dgraph](https://dgraph.io/) database, so that they can later be consulted, through the endpoints, in turn, it is responsible for transforming the structure of the nodes to a structure of a specific programming language and execute it through the [JDoodle](https://www.jdoodle.com/) api.

# Run Dgraph

Running the dgraph/standalone docker image.

Run the following command:
```
docker run --rm -it -p 8000:8000 -p 8080:8080 -p 9080:9080 dgraph/standalone:v20.11.0
```

You now have Dgraph up and running.


# JDoodle 

JDoodle is in charge of executing the program, it works for different programming languages.

### Steps to use JDoodle:

1. [Goto api](https://www.jdoodle.com/compiler-api) and sign up for the free plan. 

2. Now, copied your credentials.

3. Standing at the root of the api folder, create the .env, replacing the following with your credentials:

```
CLIENT_ID=<your_client_id>
CLIENT_SECRET=<your_client_secret>
```

Note: All this must be done after you have cloned or downloaded the repository.

# Run api

Cloned or downloaded repository.

Stand on api folder:
```
cd api
```

Run the following command:
```
go run . 
```

Now the API will run on port 10000.
```http://localhost:10000```