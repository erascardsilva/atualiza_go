package menu

import (
	"atualiza_go/variavel"
	"fmt"
	"os"
	"os/exec"
)

func Menu() {

	for {
		limpa()
		fmt.Println("▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅")
		fmt.Println("			              				  ")
		fmt.Println("			Qual seu Sistema              ")
		fmt.Println("			              				  ")
		fmt.Println("	1) Debian ou Derivado		          ")
		fmt.Println("	2) Arch ou Derivado	     			  ")
		fmt.Println("	3) Sair	            				  ")
		fmt.Println("			            				  ")
		fmt.Println("▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅")

		fmt.Scanln(&variavel.Opt)

		switch variavel.Opt {
		case 1:
			variavel.Sis = "debian"
		case 2:
			variavel.Sis = "arch"
		case 3:
			return

		}

		if variavel.Opt == 1 {
			atual_debian()
			break
		} else if variavel.Opt == 2 {
			atual_arch()
		} else {
			continue
		}

	}
}

func atual_debian() {

	for {
		limpa()
		fmt.Println("▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅")
		fmt.Println("			              				   ")
		fmt.Println("			Atualiza Sistema               ")
		fmt.Println("			              				   ")
		fmt.Println("	1) Sistema | Flatpak		           ")
		fmt.Println("	2) Sistema | Flatpak | Reset	       ")
		fmt.Println("	3) Sistema | Flatpak | Desligar	       ")
		fmt.Println("	4) Sair	            				   ")
		fmt.Println("			            				   ")
		fmt.Println("▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅")

		fmt.Scanln(&variavel.Opt2)

		switch variavel.Opt2 {
		case 1:
			executarComando("sudo apt update  -y && sudo apt upgrade  -y && flatpak update")
		case 2:
			executarComando("sudo apt update  -y && sudo apt upgrade  -y && flatpak update  && sudo reboot")
		case 3:
			executarComando("sudo apt update  -y && sudo apt upgrade  -y && flatpak update  && sudo shutdown now")
		case 4:
			return
		}
		if variavel.Opt2 >= 1 && variavel.Opt2 <= 4 {

			break
		} else {
			continue
		}

	}
}

func atual_arch() {

	for {
		limpa()
		fmt.Println("▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅")
		fmt.Println("			              				   ")
		fmt.Println("			Atualiza Sistema               ")
		fmt.Println("			              				   ")
		fmt.Println("	1) Sistema | Flatpak		           ")
		fmt.Println("	2) Sistema | Flatpak | Reset	       ")
		fmt.Println("	3) Sistema | Flatpak | Desligar	       ")
		fmt.Println("	4) Sair	            				   ")
		fmt.Println("			            				   ")
		fmt.Println("▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅")

		fmt.Scanln(&variavel.Opt2)

		switch variavel.Opt2 {
		case 1:
			executarComando("sudo pacman -Syu  --noconfirm  && flatpak update")
		case 2:
			executarComando("sudo pacman -Syu  --noconfirm  && flatpak update  && sudo reboot")
		case 3:
			executarComando("sudo pacman -Syu  --noconfirm  && flatpak update  && sudo shutdown now")
		case 4:
			return
		}
		if variavel.Opt2 >= 1 && variavel.Opt2 <= 4 {

			break
		} else {
			continue
		}

	}
}

func executarComando(comando string) {
	cmd := exec.Command("bash", "-c", comando)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Erro ao executar o comando:", err)
	}
}

func limpa() {
	limpa := exec.Command("clear")
	limpa.Stdout = os.Stdout
	limpa.Run()
}
