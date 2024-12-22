set -x

if [ "$LOCAL" = "true" ]; then
  DYNAMODB_ENDPOINT="http://localhost:4566"
else
  DYNAMODB_ENDPOINT="http://localstack:4566"
fi

echo "Creating 'users' table in dynamodb"
aws --region us-east-1 \
    --endpoint-url=$DYNAMODB_ENDPOINT dynamodb create-table \
    --table-name users \
    --attribute-definitions AttributeName=Id,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

echo "Add Luke Skywalker user"
aws --region us-east-1 \
    --endpoint-url=$DYNAMODB_ENDPOINT dynamodb put-item \
    --table-name users \
    --item '{
        "Id": {"S": "34541724-d210-47e0-9cd3-dc950344e421"},
        "FirstName": {"S": "Luke"},
        "LastName": {"S": "Skywalker"},
        "Nickname": {"S": "Starkiller"},
        "Password": {"S": "tw0M00ns"},
        "Email": {"S": "luke.skywalker@gmail.com"},
        "Country": {"S": "Tattooine"},
        "CreatedAt": {"S": "2024-07-15T07:25:55.32Z"}
    }'

echo "Add Darth Vader user"
aws --region us-east-1 \
    --endpoint-url=$DYNAMODB_ENDPOINT dynamodb put-item \
    --table-name users \
    --item '{
        "Id": {"S": "34541724-d210-47e0-9cd3-dc950344e422"},
        "FirstName": {"S": "Darth"},
        "LastName": {"S": "Vader"},
        "Nickname": {"S": "Anakin"},
        "Password": {"S": "iHateSand"},
        "Email": {"S": "darth.vader@deathstar.com"},
        "Country": {"S": "Deathstar"},
        "CreatedAt": {"S": "1977-07-15T07:25:55.32Z"}
    }'

echo "Add Storm Trooper 1"
aws --region us-east-1 \
    --endpoint-url=$DYNAMODB_ENDPOINT dynamodb put-item \
    --table-name users \
    --item '{
        "Id": {"S": "34541724-d210-47e0-9cd3-dc950344e473"},
        "FirstName": {"S": "Storm"},
        "LastName": {"S": "Trooper 1"},
        "Nickname": {"S": "Stormie"},
        "Password": {"S": "trooper"},
        "Email": {"S": "storm.trooper1@deathstar.com"},
        "Country": {"S": "Deathstar"},
        "CreatedAt": {"S": "1977-07-15T07:25:55.32Z"}
    }'

echo "Add Storm Trooper 2"
aws --region us-east-1 \
    --endpoint-url=$DYNAMODB_ENDPOINT dynamodb put-item \
    --table-name users \
    --item '{
        "Id": {"S": "34541724-d210-47e0-9cd3-dc950344e444"},
        "FirstName": {"S": "Storm"},
        "LastName": {"S": "Trooper 2"},
        "Nickname": {"S": "Stormie"},
        "Password": {"S": "trooper"},
        "Email": {"S": "storm.trooper2@deathstar.com"},
        "Country": {"S": "Deathstar"},
        "CreatedAt": {"S": "1977-07-15T07:25:55.32Z"}
    }'

