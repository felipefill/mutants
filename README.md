# Mutants

Mutants is a serverless GO app that can identify a given DNA as mutant.

This project was written in [GO][https://golang.org/] and uses [AWS Lambda][https://aws.amazon.com/lambda/] to serve its endpoints. 
The dabase is also hosted by Amazon ([RDS][https://aws.amazon.com/rds/]), there's a script under `database` folder that describes the model.

I decided to use the [serverless][https://serverless.com/] framework in order to facilitate and speed up development.

## Setup

### Dependencies

- GO
- Dep (go dependency manager)
- Serverless
- awscli (configured)

### Installing and configuring

You can usually install it by using a package manager, e.g.:
```
brew install go dep serverless awscli
```

You will have to clone this project inside a proper set up [GOPATH][https://golang.org/doc/code.html#GOPATH].

As for `awscli`, you must have it configured with credentials (that have access to Lambda related stuff):
```
aws configure
```

Last but not least, you will need to write the database info to a `serverless.env.yml` file. There's a sample included in this repo.

You can build, test and deploy using [make][https://en.wikipedia.org/wiki/Make_(software)]:

```
make build # Builds the project
make test # Run all the tests and shows code coverage
make deploy # Deploys to AWS Lambda
```

