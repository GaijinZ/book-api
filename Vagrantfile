# -*- mode: ruby -*-
# vi: set ft=ruby :

$setup = <<-SCRIPT
sudo apt-get update

wget https://dl.google.com/go/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.21.0.linux-amd64.tar.gz

/bookapi/init-scripts/env-vars.sh
source ~/.bash_profile

sudo apt-get install postgresql postgresql-contrib -y

sudo apt install make

mkdir ~/usersdb_data
mkdir ~/booksdb_data
mkdir ~/shops_books_db_data

sudo pg_createcluster --datadir=~/usersdb_data 12 usersdb
sudo pg_createcluster --datadir=~/booksdb_data 12 booksdb
sudo pg_createcluster --datadir=~/shops_books_db_data  12 shops_books_db_data

sudo pg_ctlcluster 12 usersdb start
sudo pg_ctlcluster 12 booksdb start
sudo pg_ctlcluster 12 shops_books_db_data  start

sudo -u postgres psql -p 5433 < /bookapi/init-scripts/users/postgres-init-users.sql
sudo -u postgres psql -p 5434 < /bookapi/init-scripts/books/postgres-init-books.sql
sudo -u postgres psql -p 5435 < /bookapi/init-scripts/shops/postgres-init.sql

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
apt  install golang-goprotobuf-dev

cd /bookapi
go mod tidy
SCRIPT

class VagrantPlugins::ProviderVirtualBox::Action::Network
    def dhcp_server_matches_config?(dhcp_server, config)
      true
    end
end

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/focal64"
  config.vm.synced_folder ".", "/bookapi", create: true

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "8192"
    vb.cpus = 2
  end
  
  config.vm.boot_timeout = 150
  config.vm.hostname = "ubuntu-vm"

  config.vm.network "private_network", ip: "192.168.33.10"

  config.vm.provision "docker"

  config.vm.provision "shell", inline: $setup
end
