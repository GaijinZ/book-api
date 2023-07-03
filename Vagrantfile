# -*- mode: ruby -*-
# vi: set ft=ruby :

$setup = <<-SCRIPT
sudo apt-get update

sudo snap install go --classic

bash /bookapi/init-scripts/env-vars.sh

sudo apt-get install postgresql postgresql-contrib -y

sudo -u postgres psql < /bookapi/init-scripts/users/postgres-init.sql

sudo apt install make
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

  config.vm.hostname = "ubuntu-vm"

  config.vm.network "private_network", ip: "192.168.33.10"

  config.vm.provision "docker"

  config.vm.provision "shell", inline: $setup
end