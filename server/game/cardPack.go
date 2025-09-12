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


//Contadores globais normais 
var globalNormalCounters = map[string]int{
	"AGUA": 1,
	"TERRA": 1,
	"FOGO": 1,
	"AR": 1,
	"MATO": 1,
}

//Contadores globais especias 
var globalSpecialCounters = map[string]int{
	"AGUA": 1,
	"TERRA": 1,
	"FOGO": 1,
	"AR": 1,
}

var globalDeck []string
var deckGlobalMutex sync.Mutex


func rebuildDeckUnsafe() {
	globalDeck = globalDeck[:0]
	
	//Adiciona cartas normais
	for i := 0; i < 80; i++ {
        card := normalCards[rand.Intn(len(normalCards))]
        num := globalNormalCounters[card]
        uniqueCard := fmt.Sprintf("%s #%d", card, num)
        globalDeck = append(globalDeck, uniqueCard)
		
        globalNormalCounters[card]++
    }
	
	//Adiciona cartas especiais
	for _, card := range specialCards {
		prefix := strings.Split(card, " ")[0]
		num := globalSpecialCounters[prefix]
		
		uniqueCard := fmt.Sprintf("%s #%d", card, num)
		globalDeck = append(globalDeck, uniqueCard)
		
		globalSpecialCounters[prefix]++
	}

	//Embaralha
	rand.Shuffle(len(globalDeck), func(i, j int) { 
		globalDeck[i], globalDeck[j] = globalDeck[j], globalDeck[i] 
	})

	fmt.Printf("Deck global reconstruído: %d cartas (80 normais + 20 especiais)\n", len(globalDeck))
}

func BuildGlobalDeck() {
	deckGlobalMutex.Lock()
	defer deckGlobalMutex.Unlock()
	rebuildDeckUnsafe()
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
		fmt.Println("Pouca carta... reconstruindo deck...")
		rebuildDeckUnsafe()
	}

	fmt.Println("Pacote aberto com sucesso")
	return selectedCards, nil
}

func ShowCardsGlobalDeck(){
	deckGlobalMutex.Lock()
	defer deckGlobalMutex.Unlock()
	fmt.Println("Deck global: ", globalDeck)
}

