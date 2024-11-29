
mkdir -p release/linux
mkdir -p release/windows

# Windows
fyne-cross windows -arch=amd64 -app-id="com.example.dmvp2p"
rm release/windows/dmvp2p.exe
cp fyne-cross/bin/windows-amd64/dmvp2p.exe release/windows/
zip -r release/windows/dmvp2p-windows-amd64.zip release/windows/

# Linux
go build -o ./tmp/dmvp2p main.go
rm release/linux/dmvp2p
cp tmp/dmvp2p release/linux/
sleep 1
tar -czvf release/linux/dmvp2p-linux-amd64.tar.gz release/linux


# SHA256sums + PGP
rm release/sums.txt
rm release/sums.txt.asc
sha256sum release/windows/dmvp2p-windows-amd64.zip | awk '{print $1 "  " "dmvp2p-windows-amd64.zip"}' >> release/sums.txt
sha256sum release/linux/dmvp2p-linux-amd64.tar.gz | awk '{print $1 "  " "dmvp2p-linux-amd64.tar.gz"}' >> release/sums.txt


gpg --armor --output release/sums.txt.asc --detach-sig release/sums.txt

echo "done"