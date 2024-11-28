package droplets

import (
	"fmt"
	"log"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

const PageSize = 4

type model struct {
	choices     []DropletBasicInfo // droplets in the page
	cursor      int                // which droplet our cursor is pointing at
	selected    int                // which to-do items are
	currentPage int
	lastPage    bool
	fetching    bool
}

func getDroplets(page int, pageSize int) ([]DropletBasicInfo, bool, error) {
	// This function should return a list of droplets, a boolean indicating if it is the last page, and an error
	droplets, lastPage, err := getDropletBasicInfo(page, pageSize)
	if err != nil {
		return nil, true, err
	}

	return droplets, lastPage, nil
}

func initialModel() model {
	currentPage := 1

	// droplets, lastPage, err := getDroplets(currentPage, pageSize)
	/* if err != nil {
		log.Fatalf("Failed to retrieve droplets from DigitalOcean API: %v", err)
	} */
	return model{
		choices:     []DropletBasicInfo{},
		cursor:      0,
		selected:    0,
		currentPage: currentPage,
		lastPage:    true,
		fetching:    true,
	}
}

func (m model) Init() tea.Cmd {

	return getDropletFirstPage
}

type UpdateDropletListMsg struct {
	droplets []DropletBasicInfo
	lastPage bool
}

func getDropletFirstPage() tea.Msg {
	currentPage := 1
	pageSize := PageSize
	droplets, lastPage, err := getDroplets(currentPage, pageSize)
	if err != nil {
		log.Fatalf("Failed to retrieve droplets from DigitalOcean API: %v", err)
	}

	return UpdateDropletListMsg{droplets: droplets, lastPage: lastPage}
}

func createSshTerminal(ip string){
	// private key for auth
	key, err := os.ReadFile("/home/alejandro/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}
	
	// create signer for auth 
	signer, err := ssh.ParsePrivateKey(key)
	// define auth method  
	auth := ssh.PublicKeys(signer)

	hostKeyCb := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	
	// client config
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
		auth,
		},
		HostKeyCallback: hostKeyCb,
	}


	// dial SSH connection
	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial SSH connection: %v", err)
	}
	defer conn.Close()

	// create new SSH session
	session, err := conn.NewSession()
	defer session.Close()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}
	
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to get terminal size: %v", err)
	}
	if err := session.RequestPty("xterm", width, height, ssh.TerminalModes{}); err != nil {
        log.Fatalf("Failed to request PTY: %v", err)
    }
	// Start shell
    session.Stdout = os.Stdout
    session.Stderr = os.Stderr
    session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
        log.Fatalf("Failed to start shell: %v", err)
    }
	// Wait for the session to end
    if err := session.Wait(); err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok && exitError.ExitStatus() == 130 {
            // Ignore exit status 130 (terminated by SIGINT)
            fmt.Println("SSH session finished")
        } else {
            log.Fatalf("Failed to wait for session: %v", err)
        }
    }

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateDropletListMsg:
		m.choices = msg.droplets
		m.lastPage = msg.lastPage
		m.fetching = false

		// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				return m, nil
			}

			//if we're not in the first page get the previous one
			if m.currentPage > 1 {
				m.currentPage--
				droplets, lastPage, err := getDroplets(m.currentPage, PageSize)
				if err != nil {
					log.Fatalf("Failed to retrieve droplets from DigitalOcean API: %v", err)
				}
				m.choices = droplets
				m.lastPage = lastPage
				m.cursor = len(m.choices) - 1
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
				return m, nil
			}
			//if we're not in the last page get the next one
			if m.lastPage == false {
				m.currentPage++
				droplets, lastPage, err := getDroplets(m.currentPage, PageSize)
				if err != nil {
					log.Fatalf("Failed to retrieve droplets from DigitalOcean API: %v", err)
				}
				m.choices = droplets
				m.lastPage = lastPage
				m.cursor = 0
			}

		case "enter":
			m.selected = m.cursor
			sshIntoIP := m.choices[m.selected].PublicAddr

			fmt.Println("Logging into droplet ", m.choices[m.selected].Name, "... ")

			createSshTerminal(sshIntoIP)
			
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// Header
	s := "What Droplet you want to log into?\n\n"
	if m.fetching {
		s += "Please wait while we fetch the droplets...\n"
	} else {
		if m.currentPage > 1 && m.cursor == 0 {
			s += "\nPress ⬆️ up to load previous droplets.\n"
		} else {
			s += "\n \n"
		}
	}
	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s [%s/%s] => (Region: %s, ID: %d)\n", cursor, choice.Name, choice.PrivateAddr, choice.PublicAddr, choice.Region, choice.ID)
	}

	// The footer
	if m.lastPage == false && m.cursor == len(m.choices)-1 {
		s += "\nPress ⬇️ down to load more droplets.\n"
	} else {
		s += "\n \n"
	}

	s += "\nPress ⬆️ up or ⬇️ down arrows to select, enter to ssh into, or q to quit.\n"

	// Send the UI for rendering
	return s
}

func SshDropletUi() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("There's been an error logging into the droplet: %v", err)
		os.Exit(1)
	}
}
