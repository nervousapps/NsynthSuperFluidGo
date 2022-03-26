setup_audio() {
    if grep -q '^dtparam=audio=on' /boot/config.txt
    then
        echo 'Disabling on board audio'
        sed -i -e 's/dtparam=audio=on/#dtparam=audio=on/' /boot/config.txt
    fi

    if ! grep -q "^dtoverlay=iqaudio-dac" /boot/config.txt
    then
        echo 'Enabling I2S + I2C audio'
        printf "\n# Enable I2S + I2C audio\ndtoverlay=iqaudio-dac\n" >> /boot/config.txt
    fi
}


setup_i2c() {
    if grep -q '^#dtparam=i2c_arm=on' /boot/config.txt
    then
        echo 'Enabling I2C'
        sed -i -e 's/#dtparam=i2c_arm=on/dtparam=i2c_arm=on/' /boot/config.txt
    fi

    if ! grep -q '^dtparam=i2c_arm_baudrate=640000' /boot/config.txt
    then
        echo 'Setting I2C to high speed'
        printf "\n# Set I2C to high speed\ndtparam=i2c_arm_baudrate=640000\n" >> /boot/config.txt
    fi

    if ! grep -q '^i2c-dev' /etc/modules
    then
        echo 'Enabling i2c-dev module'
        printf "i2c-dev\n" >> /etc/modules
    fi
}


setup_serial() {
    if grep -q 'console=serial' /boot/cmdline.txt
    then
        echo 'Disabling console serial'
        sed -i -e 's/console=serial[^ ]* //' /boot/cmdline.txt
    fi

    if ! grep -q '^dtoverlay=pi3-miniuart-bt' /boot/config.txt
    then
        echo 'Enabling serial MIDI'
        printf "\n# Enable seiral MIDI\nenable_uart=1\ndtoverlay=pi3-miniuart-bt\ndtoverlay=midi-uart0\n" >> /boot/config.txt
    fi
}

install_all_deps() {
  sudo apt-get update
  sudo apt-get upgrade -y
  sudo apt-get install -y pulseaudio libpulse-dev libportmidi-dev osspd alsa-utils alsa-oss alsa-tools jack jackd2 jack-tools a2jmidid libjack-jackd2-dev libinstpatch-1.0 libglib2.0-dev
  sudo apt install libasound2-dev
  sudo apt install git check libglib2.0-dev libreadline-dev libudev-dev libsystemd-dev libusb-dev cmake build-essential libsndfile-dev swami
  sudo apt-get autoremove -y

  sudo gpasswd -a $USER audio
}

setup_service() {
  sudo cp nsynthsuperhard.service /lib/systemd/system/nsynthsuperhard.service
  sudo chmod 644 /lib/systemd/system/nsynthsuperhard.service
  sudo systemctl daemon-reload

  sudo pigpiod
  sudo systemctl enable nsynthsuperhard.service
}

build_fluidsynth() {
    sudo apt-get install -y python-pip
    pip install lastversion
    FLUIDSYNTH_VERSION=$(lastversion https://github.com/FluidSynth/fluidsynth)
    wget https://github.com/FluidSynth/fluidsynth/archive/refs/tags/v${FLUIDSYNTH_VERSION}.tar.gz
    tar -xvf v${FLUIDSYNTH_VERSION}.tar.gz
    cd fluidsynth-${FLUIDSYNTH_VERSION}
    mkdir build
    cd build
    cmake ..
    make
    sudo make install
    make check
    
    cd ../../
    rm -rf fluidsynth-${FLUIDSYNTH_VERSION}
    rm -rf v2.2.6.tar.gz
}

# install_all_deps
# setup_i2c
# setup_audio
# setup_serial

build_fluidsynth

# setup_service

# sudo reboot
