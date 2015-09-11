# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  numNodes = ENV["NODES"].to_i
  puts numNodes
  config.vm.define :master do |master|
    master.vm.box = "ubuntu/trusty64"
    master.vm.network :private_network, ip: "192.168.1.100"
    master.vm.hostname = "master"

    master.vm.provision "docker"
    master.vm.provision :shell, :path => "all.sh"
    master.vm.provision :shell, :path => "hosts.sh", :args => "'%d'" % numNodes
    master.vm.provision :shell, :path => "master.sh", :args => "'%d'" % numNodes
  end

  1.upto(numNodes) do |num|
    nodeName = ("slave" + num.to_s).to_sym
    val = num + 100
    config.vm.define nodeName do |node|
      node.vm.box = "ubuntu/trusty64"
      node.vm.network :private_network, ip: "192.168.1." + val.to_s
      node.vm.hostname = "slave" + num.to_s

      node.vm.provision "docker"
      node.vm.provision :shell, :path => "all.sh"
      node.vm.provision :shell, :path => "hosts.sh", :args => "'%d'" % numNodes
      node.vm.provision :shell, :path => "slave.sh", :args => "'%d'" % num
    end
  end
end
