# -*- mode: ruby -*-
# vi: set ft=ruby :
required_plugins = %w(vagrant-vbguest)
required_plugins.each do |plugin|
  system "vagrant plugin install #{plugin}" unless Vagrant.has_plugin? plugin
end

Vagrant.configure("2") do |config|

  config.vm.define "keycloak" do |d|
    d.vm.box = "centos/7"
    d.ssh.insert_key = false
    d.vm.box_version = "1905.1"
    d.vm.hostname = "keycloak"
    d.vm.network "private_network", ip: "10.100.196.60"

    d.vm.provision :shell, path: "scripts/passwordAuthentication.sh"
    d.vm.provision :shell, path: "scripts/bootstrap_misc.sh"

    d.vm.provider "virtualbox" do |v|
      v.memory = 4096
      v.cpus = 4
    end
  end

  
  if Vagrant.has_plugin?("vagrant-vbguest")
    config.vbguest.auto_update = true
    config.vbguest.no_remote = true
  end

end
