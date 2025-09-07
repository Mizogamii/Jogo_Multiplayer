package game

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

var normalCards = []string{"AGUA", "TERRA", "FOGO", "AR", "MATO"}

var specialCards = []string{
	//Água
	"AGUA BENTA", "AGUA LIMPA", "AGUA GLACIAL", "AGUA DOCE", "AGUA FORTE",

	//Terra
	"TERRA DIAMENTE", "TERRA VULCANICA", "TERRA FERTIL", "TERRA VERMELHA", "TERRA FORTE",

	//Fogo
	"FOGO SAGRADO", "FOGO HADES", "FOGO ETERNO", "FOGO DRAGAO", "FOGO FORTE",

	//Ar
	"AR FEDENDO", "AR POLUIDO", "AR LIMPO", "AR TORNADO", "AR FORTE",

}

var specialCounters = map[string]int{
	"AGUA": 1,
	"TERRA": 1,
	"FOGO": 1,
	"AR": 1,
}

var globalDeck []string
var deckGlobalMutex sync.Mutex

func BuildGlobalDeck() {
	deckGlobalMutex.Lock()
	defer deckGlobalMutex.Unlock()
	for i := 0; i < 80; i++ {
		card := normalCards[rand.Intn(len(normalCards))]
		globalDeck = append(globalDeck, card)
	}
	
	for _, card := range specialCards {
		prefix := strings.Split(card, " ")[0]
		num := specialCounters[prefix]
		
		uniqueCard := fmt.Sprintf("%s #%d", card, num)
		globalDeck = append(globalDeck, uniqueCard)
		
		specialCounters[prefix]++
	}

	//Embaralha
	rand.Shuffle(len(globalDeck), func(i, j int) { 
		globalDeck[i], globalDeck[j] = globalDeck[j], globalDeck[i] 
	})

	fmt.Printf("Deck global criado: %d cartas (80 normais + 20 especiais)\n", len(globalDeck))
}

func OpenPack(playerName string)([]string, error){
	deckGlobalMutex.Lock()
	defer deckGlobalMutex.Unlock()
	
	if len(globalDeck) < 3{
		return nil, fmt.Errorf("acabou as cartas")
	}
	selectedCards := make([]string, 3)
	copy(selectedCards, globalDeck[:3])

	globalDeck = globalDeck[3:]

	if len(globalDeck) < 10{
		fmt.Println("Pouca carta...")
		BuildGlobalDeck()
	}

	fmt.Println("Pacote aberto com sucesso")
	return selectedCards, nil
}