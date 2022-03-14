#!/usr/bin/python3
# -*- coding: utf-8 -*-
import os
import sys
import shutil

# global vars/filepaths
bin_dir = '/usr/bin/'
config_dir = '$HOME/.config/'
required_files = ['wpe.py', 'wpe.json']
gif_file = 'include/foo-Wallpaper-Feh-Gif/back4.sh'


def cmd_parser():
    commands = {    # available commands, determin the installation process
        'gif': False,
        'no-config': False
    }

    if len(sys.argv) <= 1:
        return commands

    start = True
    for arg in sys.argv:
        if start:
            start = False
            continue
        if arg == '--gif':
            commands['gif'] = True
        elif arg == '--no-config':
            commands['no-config'] = True
        else:
            print(f'{arg} is not an existing command')

    return commands


def make_executable(file):
    try:
        st = os.stat(file)
        os.chmod(file, st.st_mode | stat.S_IEXEC)
    except Exception as ex:
        exit(f'Encountered an error while changing file permissions.\nError: {ex}')


def install(config):
    make_executable(required_files[0])
    if config['gif']:
        if not os.path.exists(gif_file):
            #TODO download the script instead of throwing an error.
            exit(f'Script for gif execution not found {gif_file}.\nPlease make sure to clone the repository recursivly for gif support.')
        make_executable(gif_file)
        print('Moving gif script...')
        shutil.copyfile(gif_file, bin_dir)
    print('Moving wpe script...')
    shutil.copyfile(required_files[0], bin_dir)
    if not config['no-config']:
        print('Moving wpe config...')
        shutil.copyfile(required_files[1], config_dir)


# init :: make sure everything is set so we can install wpe without any problems
def init():
    if os.geteuid() != 0:   # root?
        exit("You need to have root privileges to run this script.\nPlease try again, this time using 'sudo'. Exiting.")

    for file in required_files:
        if not os.path.exists(file):
            exit(f'File {file} does not exist.\nPlease make sure you are in the same directory as the installation files or freshly clone WallpaperEngine from github.')


def main():
    config = cmd_parser()
    install(config)
    print('Done. WallpaperEngine(.py) is now available as wpe')
    exit(0)


if __name__ == '__main__':
    init()
    main()
