# Enter build directory
cd build || exit 1

# Cleanup possible old result
rm -f nitori.img

# Allocate image
fallocate -l $(($(du -b init | cut -f1) + 67108864)) nitori.img

# Find available loop device
LOOP=$(sudo losetup -f)

# Setup partitions
sudo parted nitori.img -- mktable gpt
sudo parted nitori.img -- mkpart efi fat32 1MiB 100%
sudo parted nitori.img -- set 1 esp on

# Attach loop device
sudo losetup --partscan "$LOOP" nitori.img

# Make EFI system partition
sudo mkfs.vfat -F 32 -n EFI "$LOOP"p1 || exit 1

# Mount the ESP
mkdir "mount"
sudo mount -t vfat "$LOOP"p1 mount || exit 1

# Deploy files
sudo cp -r ../assets/os/boot mount/
sudo cp -r ../assets/os/efi mount/
sudo cp ../../../linux/arch/x86/boot/bzImage mount/vmlinuz
sudo cp init mount/
sudo mkdir mount/dev

# Unmount and cleanup
sudo umount mount
sudo rmdir mount
sudo losetup -d "$LOOP"
