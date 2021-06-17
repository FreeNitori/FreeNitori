# Enter build directory
cd build || exit 1

# ESP size
ESP_SIZE=16777216

# Cleanup possible old result
rm -f nitori.img

# Prepare root
mkdir -p root/{dev,bin,sbin,var,proc}
cp init root/sbin/
cp freenitori nitorictl root/bin/

# Build rootfs image
fakeroot mksquashfs root nitori.sqsh
rm -r root

# Allocate image
fallocate -l $(($(du -b nitori.sqsh | cut -f1) + "$ESP_SIZE" + 20480)) nitori.img

# Setup partitions
parted nitori.img -- mktable gpt
parted nitori.img -- mkpart efi fat32 4MiB "$((ESP_SIZE - 1))"B
parted nitori.img -- mkpart linux ext4 "$ESP_SIZE"B 100%
parted nitori.img -- set 1 esp on
sgdisk --partition-guid=2:726c8d61-f9e6-429a-86b6-d33f774e91b4 nitori.img

# Find available loop device
LOOP=$(sudo losetup -f)

# Attach loop device
sudo losetup --partscan "$LOOP" nitori.img

# Make EFI system partition
sudo mkfs.vfat -F 16 -n EFI "$LOOP"p1 || exit 1

# Mount the ESP
mkdir "mount"
sudo mount -t vfat "$LOOP"p1 mount || exit 1

# Prepare ESP
sudo cp -r ../assets/os/esp/* mount/

# Unmount and cleanup
sudo umount mount
sudo rmdir mount

# Copy rootfs into partition
sudo sh -c "cat nitori.sqsh > ${LOOP}p2"
rm nitori.sqsh

# Detach the loop device
sudo losetup -d "$LOOP"
