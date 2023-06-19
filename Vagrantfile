# -*- mode: ruby -*-
# vi: set ft=ruby :

$setup = <<-SCRIPT

sudo apt update
sudo apt-get install postgresql postgresql-contrib -y

sudo -u postgres psql < /bookapi/init-scripts/users/postgres-init.sql

SCRIPT

class VagrantPlugins::ProviderVirtualBox::Action::Network
    def dhcp_server_matches_config?(dhcp_server, config)
      true
    end
end

  Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/xenial64"
    config.vm.synced_folder ".", "/bookapi", create: true
    
    config.vm.define "ubuntu_vm" do |ubuntu|
      ubuntu.vm.hostname = "ubuntu-vm"
      ubuntu.vm.network "private_network", ip: "192.168.33.10"
      
      ubuntu.vm.provider "virtualbox" do |vb|
        vb.memory = "8192"
        vb.cpus = 2
      end
    end
    config.vm.provision "shell", inline: $setup
  end