import subprocess


def run(cmd):
    print("Running:", cmd)
    subprocess.run(cmd, shell=True, check=True)


def main():
    run("multipass stop --all")

    run("multipass delete --all")

    run("multipass purge")

    print("Finished Purge")


if __name__ == "__main__":
    main()