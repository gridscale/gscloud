class Gscloud < Formula
  desc "Official command-line interface for the gridscale API"
  homepage "https://gridscale.io/"
  url "https://github.com/gridscale/gscloud/archive/refs/tags/v0.12.0.tar.gz"
  sha256 "20927acda1fff7372bd6de11dcd40b0b6143aa6668d88b79d181cd9ccf5440f4"
  license "MIT"
  head "https://github.com/gridscale/gscloud.git", branch: "master"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/gridscale/gscloud/cmd.Version=#{version}
    ]

    system "go", "build", *std_go_args(ldflags: ldflags)

    # currently gscloud needs at least an empty config to run
    touch "config.yaml"
    (zsh_completion/"_gscloud").write `#{bin}/gscloud completion zsh`
    (bash_completion/"gscloud").write `#{bin}/gscloud completion bash`
  end

  test do
    # currently gscloud needs at least an empty config to run
    touch "config.yaml"
    assert_match "Version:\t#{version}", shell_output("#{bin}/gscloud version")
    assert_match "gscloud lets you manage", shell_output("#{bin}/gscloud help")
  end
end
