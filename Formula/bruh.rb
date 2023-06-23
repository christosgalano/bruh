# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Bruh < Formula
  desc "Command-line tool for scanning and updating the API version of Azure resources in Bicep files"
  homepage "https://github.com/christosgalano/bruh"
  version "0.1.0"
  license "Apache 2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/christosgalano/bruh/releases/download/v0.1.0/bruh_darwin_amd64.tar.gz"
      sha256 "b786b6651b69cf0eb3ab8ea0fac26013c5ef441a964e84ac1f6f5937e2e3e3fd"

      def install
        bin.install "bruh"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/christosgalano/bruh/releases/download/v0.1.0/bruh_darwin_arm64.tar.gz"
      sha256 "ee7a3e04fbe74a524d9bf57814c278d19744452b622fbc01010c15348f03f7fc"

      def install
        bin.install "bruh"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/christosgalano/bruh/releases/download/v0.1.0/bruh_linux_arm64.tar.gz"
      sha256 "5d5a9465e9fd135bbbaa2f91e6f5ba30ad476fd2a3000da107e0c47d54fac625"

      def install
        bin.install "bruh"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/christosgalano/bruh/releases/download/v0.1.0/bruh_linux_amd64.tar.gz"
      sha256 "4fb110bfe37b75f1018c7134b4be0652d5e20edec8b637be2cb87b3188bdc35c"

      def install
        bin.install "bruh"
      end
    end
    if Hardware::CPU.arm? && !Hardware::CPU.is_64_bit?
      url "https://github.com/christosgalano/bruh/releases/download/v0.1.0/bruh_linux_arm.tar.gz"
      sha256 "d7c68fd252fec843197942b78aadd5595d723741b4c3220c212f29bedf6ca428"

      def install
        bin.install "bruh"
      end
    end
  end

  test do
    system "#{bin}/bruh", "--help"
  end
end