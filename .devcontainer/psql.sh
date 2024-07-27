# check if postgres is installed
if ! command -v psql &> /dev/null
then
    echo "Postgres is not installed. Installing..."
    sudo apt-get update
    sudo apt-get install -y postgresql-client
    echo "Postgres client installed."
else
    echo "Postgres is already installed."
fi

