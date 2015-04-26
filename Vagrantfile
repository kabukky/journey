# -*- mode: ruby -*-
# vi: set ft=ruby ts=2 sw=2 tw=0 et :

journey_archive = ENV['JOURNEY_ARCHIVE'] ||= "journey-linux-debian-wheezy-stable-amd64.zip"
journey_version = ENV['JOURNEY_VERSION'] ||= "0.1.5"

script = <<EOF
set -e
apt-get -q update && apt-get -y -q install unzip
cd /home/vagrant
wget --no-clobber --no-verbose https://github.com/kabukky/journey/releases/download/v${2}/${1}
sudo -u vagrant unzip -o -q ${1}
sed -i 's;"Url":"http://127.0.0.1:8084";"Url":"http://0.0.0.0:8084";' ${1%%.zip}/config.json
sudo -u vagrant bash -c "${1%%.zip}/journey &"
EOF

Vagrant.configure(2) do |config|
  config.vm.box = "opscode-debian-7.8.0"
  config.vm.box_url = "https://opscode-vm-bento.s3.amazonaws.com/vagrant/virtualbox/opscode_debian-7.8_chef-provisionerless.box"
  config.vm.hostname = "journey"

  config.vm.provider "virtualbox" do |vbox|
    vbox.customize ["modifyvm", :id, "--cpuexecutioncap", "50"]
    vbox.customize ["modifyvm", :id, "--memory", "128"]
  end

  config.vm.network :private_network, ip: "10.0.20.2"

  config.vm.provision "shell", inline: script, args: "#{journey_archive} #{journey_version}"
end

