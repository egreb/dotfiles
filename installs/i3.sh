set -e

# i3 is an X11 window manager. On a fresh Ubuntu (26.04+ ships Wayland-only
# GNOME with no standalone Xorg server), you MUST install a real Xorg server
# or the i3 session never appears in the GDM login gear/cog menu.
# See README.md in this dir for the full story.

echo "installing Xorg (required for i3, not installed by default on Wayland-only Ubuntu)"
sudo apt install -y xserver-xorg xinit

echo "installing i3"
# i3 is in Ubuntu's universe repo. The sur5r repo (per i3wm.org docs) often
# lacks builds for brand-new Ubuntu releases, so prefer the distro package.
sudo apt install -y i3

echo "done. reboot, then pick i3 from the gear menu on the GDM login screen."
