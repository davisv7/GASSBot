package main

import (
	"bufio"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Individual struct {
	given          map[[2]int]int
	board          [][]int
	mrate, fitness int
}

func (me *Individual) initialize() {
	me.fill()
	me.mutate()
	me.getFitness()
	//fmt.Println(me.board)
}

func (me *Individual) fill() {

	for i, row := range me.board {
		for j := range row {
			if me.board[i][j] == 0 {
				me.board[i][j] = rand.Intn(9) + 1
			}
		}
	}
}

func (me *Individual) mutate() {
	for i, row := range me.board {
		for j := range row {
			roll := rand.Intn(100)
			_, ok := me.given[[2]int{i, j}]
			if roll < me.mrate && !ok {
				me.board[i][j] = rand.Intn(9) + 1
			}
		}
	}
}

func (me *Individual) getFitness() {
	score := 0
	set := make(map[int]int)
	for _, row := range me.board {
		for _, col := range row {
			set[col] = 1
		}
		score = score + 2*(9-len(set))
		set = make(map[int]int)
	}
	for i := range me.board {
		for j := range me.board[i] {
			set[me.board[j][i]] = 1
		}
		score = score + 2*(9-len(set))
		set = make(map[int]int)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				for l := 0; l < 3; l++ {
					set[me.board[l+3*j][k+3*i]] = 1
					//fmt.Println(l+3*j,k+3*i)
				}
			}
			score = score + (9 - len(set))
			set = make(map[int]int)
		}
	}
	//fmt.Println(score)
	me.fitness = score

}

func (p *Population) mutateAll() {
	for i := 0; i < p.popsize; i++ {
		p.population[i].mutate()
	}
}

type Population struct {
	popsize, gens, parentsize, mrate, generation int
	given                                        map[[2]int]int
	solution                                     [][]int
	problem                                      [][]int
	population                                   []Individual
	top_individuals                              []Individual
}

func (p *Population) importProblem() {
	file, err := os.Open("easyproblem.in")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // internally, it advances token based on sperator
		str := strings.Split(scanner.Text(), "")
		row := make([]int, 0)
		for i := range str {
			s, err := strconv.Atoi(str[i])
			if err != nil {
				row = append(row, 0)
			} else {
				row = append(row, s)
			}
		}
		p.problem = append(p.problem, row) // token in unicode-char
	}
}

func (p *Population) importSolution() {
	file, err := os.Open("easysolution.in")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // internally, it advances token based on seperator
		str := strings.Split(scanner.Text(), "")
		row := make([]int, 0)
		for i := range str {
			s, _ := strconv.Atoi(str[i])
			row = append(row, s)
		}
		p.solution = append(p.solution, row) // token in unicode-char
	}
}

func copyBoard(board [][]int) [][]int {
	copyboard := make([][]int, 9)
	for i := range board {
		copyboard[i] = make([]int, len(board[i]))
		copy(copyboard[i], board[i])
	}
	return copyboard
}

func (p *Population) printBoard(board [][]int) {
	for i := range board {
		row := make([]int, 0)
		for j := range board[i] {
			if board[i][j] == p.solution[i][j] {
				row = append(row, board[i][j])
			} else {
				row = append(row, 0)
			}
		}
		fmt.Println(row)
	}
}

func copyIndividual(ind Individual) Individual {
	return Individual{given: ind.given, mrate: ind.mrate, board: copyBoard(ind.board), fitness: ind.fitness}
}

func (p *Population) populate() {
	for i := 0; i < p.popsize; i++ {
		ind := Individual{given: p.given, mrate: p.mrate, board: copyBoard(p.problem)}
		ind.initialize()
		p.population = append(p.population, ind)
	}
}

func (p *Population) repopulate() {
	for i := p.parentsize; i < p.popsize; i++ {
		parentOne, parentTwo := p.top_individuals[rand.Intn(p.parentsize)], p.top_individuals[rand.Intn(p.parentsize)]
		ind := Individual{given: p.given, mrate: p.mrate, board: crossoverBoard(parentOne, parentTwo)}
		ind.mutate()
		ind.getFitness()
		p.population[i] = ind
	}
}

func crossoverBoard(p1 Individual, p2 Individual) [][]int {
	board := make([][]int, 0)
	for i := range p1.board {
		row := make([]int, 0)
		for j := range p1.board[i] {
			index := rand.Intn(2)
			row = append(row, []int{p1.board[i][j], p2.board[i][j]}[index])
		}
		board = append(board, row)
	}
	return board

}

func (p *Population) findGiven() {
	for i, row := range p.problem { // i is the index, e the element
		for j, col := range row {
			if col != 0 {
				var pos = [2]int{i, j}
				p.given[pos] = 1
			}
		}
	}
}

func (p *Population) getFitnesses() {
	for i := 0; i < p.popsize; i++ {
		p.population[i].getFitness()
	}
}

type ByFitness []Individual

func (me ByFitness) Len() int           { return len(me) }
func (me ByFitness) Less(i, j int) bool { return me[i].fitness < me[j].fitness }
func (me ByFitness) Swap(i, j int)      { me[i], me[j] = me[j], me[i] }
func (p *Population) printFitnesses() {
	for i := range p.population {
		fmt.Println(p.population[i].fitness)
	}
}
func (p *Population) getTops() {
	sort.Sort(ByFitness(p.population))
	p.top_individuals = make([]Individual, 0)
	for i := 0; i < p.parentsize; i++ {
		p.top_individuals = append(p.top_individuals, copyIndividual(p.population[i]))
	}
}

func (p *Population) initialize() {
	stop := false
	stuck := 0
	prevTop := 0
	p.given = make(map[[2]int]int)
	p.importProblem()
	p.importSolution()
	p.findGiven()
	p.populate()
	for i := 0; i < p.gens; i++ {
		//p.getFitnesses()
		p.getTops()
		p.repopulate()
		if prevTop == p.top_individuals[0].fitness {
			stuck = stuck + 1
		} else {
			stuck = 0
		}
		for _, ind := range p.top_individuals {
			if ind.fitness == 0 {
				stop = true
				break
			}
		}

		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
		p.printBoard(p.top_individuals[0].board)
		fmt.Println(p.top_individuals[0].fitness, stuck)
		prevTop = p.top_individuals[0].fitness
		if stop {
			break
		}
		if stuck == 100 {
			fmt.Println("Shuffling...")
			time.Sleep(1 * time.Second)
			for i := 0; i < 10; i++ {
				p.mutateAll()
			}
			stuck = 0
		}

		//fmt.Println(p.top_individuals[0].fitness)
	}
}

func init() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("Aw fuck :(")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func main() {
	start := time.Now()
	x := Population{popsize: 10000, gens: 10000, parentsize: 1000, mrate: 4}
	//fmt.Println(x)
	x.initialize()
	elapsed := time.Since(start)
	log.Printf("GA took %s", elapsed)
	//x.population[0].board[0][0] = "69"
	//fmt.Println(x.problem)
	//fmt.Println(x.solution)
	//fmt.Println(x.given)
	//fmt.Println(x.population)
	//fmt.Println(x.top_individuals)

	//bufio.NewReader(os.Stdin).ReadBytes('\n') 
}
