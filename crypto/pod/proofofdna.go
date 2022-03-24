package pod

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/crypto/sha3"
	"golang.org/x/image/draw"
	"image"

	"image/png"
	"log"

	"os"

	"encoding/json"
	"fmt"
	"github.com/mb-14/gomarkov"
	"golang.org/x/crypto/blake2b"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"
)

/*
	this will host interface and primitives
for the temper proof certificate
*/

var mediaMutationRate = 0.008
var FitnessLimit = 1000

type TemperProofParams struct {
	MutationRate   float64
	ParseDuration  time.Duration
	PopulationSize uint
	MaxFitness     float64
}

type TemperProof interface {
	setMutationRate(rate float64)
	setParseDuration(dur time.Duration)
	setPopulationSize(size uint)
	setMaxFitness(fit float64)
	setTarget(target []byte) string
	createOrganism(target []byte) []Organism
	createPopulation(target []byte, popSize uint) (population Organism)
	crossover(d1 Organism, d2 Organism) Organism
	naturalSelection(pool []Organism, population []Organism, target []byte) []Organism
	loadDictionary()
	train() *TrainingDataSet
	getBest(population []Organism) Organism
	getMutationRate() float64
}

type TrainingDataSet []TrainingData

type TrainingData struct {
	blockText   string
	contentText string
}

type ITrainer interface {
	GetTrainingSet() map[int]string
	calculateDifficultyOfTarget(target string) (score float64)
	Train() string
	buildMarkovModel() (*gomarkov.Chain, error)
	saveMarkovModel(chain *gomarkov.Chain)
	loadModel() (*gomarkov.Chain, error)
	generateProof(chain *gomarkov.Chain) string
}

type Trainer struct {
	ITrainer
}

func (t *Trainer) GetTrainingSet() map[int]string {
	// todo change quotes.json to real data source (content and blocks)
	data, err := ioutil.ReadFile("./dict/quotes.json")
	if err != nil {
		panic(err)
	}

	mapper := make(map[int]string, len(data))
	splitted := strings.Split(string(data), "\n")

	for i := 0; i < len(splitted); i++ {

		mapper[i] = string(splitted[i])

	}
	return mapper

}

func (t *Trainer) calculateDifficultyOfTarget(target string) (score float64) {
	return 0
}

func (t *Trainer) Train() string {

	chain, err := t.buildMarkovModel()
	if err != nil {
		panic(err)
	}
	t.saveMarkovModel(chain)
	chain, err = t.loadModel()
	return t.generateProof(chain)

}

func (t *Trainer) buildMarkovModel() (chain *gomarkov.Chain, err error) {

	trainingData := t.GetTrainingSet()
	chain = gomarkov.NewChain(1)
	var wg sync.WaitGroup
	wg.Add(len(trainingData))

	for _, td := range trainingData {
		td := td
		go func() {
			defer wg.Done()
			chain.Add(strings.Split(td, " "))
		}()
	}
	wg.Wait()

	return chain, nil

}

func (t *Trainer) saveMarkovModel(chain *gomarkov.Chain) {
	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (t *Trainer) loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func (t *Trainer) generateProof(chain *gomarkov.Chain) string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}
	return fmt.Sprintf(strings.Join(tokens[1:len(tokens)-1], " "))
}

type IOrganism interface {
	calculateFitness(target []byte)
	mutate(mutationRate float64)
}

type Organism struct {
	DNA     []byte
	Fitness float64
	ParentA []byte
	ParentB []byte
	IOrganism
}

type DnaOperator struct {
	params    TemperProofParams
	processor *process
	trainer   *Trainer
}

func NewOperator(params TemperProofParams) *DnaOperator {
	op := new(DnaOperator)
	op.params = params
	op.trainer = NewTrainer()
	op.processor = new(process)
	return op
}

func NewTrainer() *Trainer {
	return &Trainer{}
}

type tProof struct {
	TemperProof
}

type operator struct {
	*DnaOperator
}

type process struct {
	TemperProof
}

func (op *process) setTarget(trainer *Trainer) string {

	target := trainer.Train()
	buf := bytes.NewBuffer(make([]byte, 32))
	i := 0
	if len(target) > 16 && len(target) <= 32 {
		tBytes := []byte(target)
		for i < 16 {
			buf.WriteByte(tBytes[i])

			i++
		}

		target = string(buf.Bytes())

	} else {
		op.setTarget(trainer)
	}
	return target

}

func (op *process) createOrganism(target []byte) (organism Organism) {
	ba := make([]byte, len(target))
	for i := 0; i < len(target); i++ {
		ba[i] = byte(rand.Intn(95) + 32)
	}

	organism = Organism{
		DNA:     ba,
		Fitness: 0,
	}

	organism.calculateFitness(target)
	return
}

func (op *process) createPopulation(target []byte, popSize uint) (population []Organism) {

	population = make([]Organism, popSize)
	for i := 0; i < int(popSize); i++ {
		population[i] = op.createOrganism(target)
	}
	return

}

func (op *process) createGenePool(maxFitness float64, population []Organism, target []byte) (pool []Organism) {
	pool = make([]Organism, 0)
	for i := 0; i < len(population); i++ {
		population[i].calculateFitness(target)
		num := int((population[i].Fitness / maxFitness) * 100)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}

	//A.genePool = pool

	return
}

func (op *process) naturalSelection(mutationRate float64, pool []Organism, population []Organism, target []byte) []Organism {
	next := make([]Organism, len(population))
	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		pa, _ := json.Marshal(a)
		pb, _ := json.Marshal(b)
		next[i].ParentA = pa
		next[i].ParentB = pb

		child := op.crossover(a, b)
		child.mutate(mutationRate)
		child.calculateFitness(target)
		next[i] = child

	}
	return next
}

func (op *process) crossover(d1 Organism, d2 Organism) Organism {
	child := Organism{
		DNA:     make([]byte, len(d1.DNA)),
		Fitness: 0,
	}

	mid := rand.Intn(len(d1.DNA))
	for i := 0; i < len(d1.DNA); i++ {
		if i > mid {
			child.DNA[i] = d1.DNA[i]
		} else {
			child.DNA[i] = d2.DNA[i]
		}
	}
	return child
}

func (op *process) getBest(population []Organism) Organism {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}

func (o *Organism) mutate(mutationRate float64) {
	for i := 0; i < len(o.DNA); i++ {
		if rand.Float64() < mutationRate {
			o.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}

func (o *Organism) calculateFitness(target []byte) {

	score := 0
	for i := 0; i < len(o.DNA); i++ {
		if o.DNA[i] == target[i] {
			score++
		}
	}
	o.Fitness = float64(score) / float64(len(o.DNA))
	return

}

func GetProof(params TemperProofParams, blockTime time.Time) (string, string) {

	opp := NewOperator(params)

	//target := []byte(opp.processor.setTarget(opp.trainer))

	//dummy target for tests only
	target := []byte("something in the way of this shit")
	startTime := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	pop := opp.processor.createPopulation(target, params.PopulationSize)
	found := false
	//generation starts at 0
	generation := 0

	var elapsed time.Duration

	for !found {
		generation++
		bo := opp.processor.getBest(pop)
		elapsed = time.Since(startTime)
		blockTime.Add(elapsed)
		fmt.Printf("\r generation: %d | %s | fitness: %2f | elapsed: %2f", generation, string(bo.DNA), bo.Fitness, elapsed)
		//dur, _ := time.ParseDuration("30s")

		if bytes.Compare(bo.DNA, target) == 0 /*|| elapsed.Seconds() <= dur.Seconds() */ {
			found = true
			oBytes, _ := json.Marshal(bo)
			pop, _ := json.Marshal(pop)
			h, _ := blake2b.New256(nil)
			h.Write(pop)
			buf := new(bytes.Buffer)
			buf.Write(h.Sum(nil))
			h.Write(oBytes)
			buf2 := new(bytes.Buffer)
			buf2.Write(h.Sum(nil))

			bufStr := buf.String()
			bufStr2 := buf2.String()
			buftotal := bufStr + bufStr2

			hmac, _ := blake2b.New256([]byte(buftotal))

			hmac.Write([]byte(buftotal))
			encoded := encodeProof(buftotal, hmac.Sum(nil))
			//spew.Dump(encoded)

			proof := NewProof(encoded)

			return proof.hash, proof.hmac

		} else {
			maxFitness := bo.Fitness
			pool := opp.processor.createGenePool(maxFitness, pop, target)
			pop = opp.processor.naturalSelection(params.MutationRate, pool, pop, target)
		}

	}

	return "", ""

}

func encodeProof(buftotal string, hmac []byte) []byte {

	var buf bytes.Buffer

	// prefix for a proof is always 0x0111
	_, err := buf.Write([]byte{0x01, 0x11})
	if err != nil {
		return nil
	}

	w1 := []byte(buftotal)
	log.Printf("bufferlen:%x", len(w1))

	_, err = buf.Write(w1)
	if err != nil {
		return nil
	}
	// hmac prefix is always 0x0222
	_, err = buf.Write([]byte{0x02, 0x22})
	if err != nil {
		return nil
	}
	_, err = buf.Write(hmac)
	if err != nil {
		return nil
	}

	buf.Cap()

	return buf.Bytes()

}

func NewProof(data []byte) *Proof {
	p := new(Proof)

	if data[0] == 0x01 && data[1] == 0x011 {
		//ok its a proof

		// we remove 0:0 and 0:1 to get the full bytes and hmac
		proofData := data[2:66]
		hmacData := data[68:]

		p.hash = hex.EncodeToString(proofData)
		p.hmac = hex.EncodeToString(hmacData)

		if p.verifyProof() {
			p.isVerified = true
			return p
		}
		p.isVerified = false
		panic("error invalid proof!")

	}

	log.Fatalf("bad proof format")
	return nil

}

func (p *Proof) verifyProof() bool {

	buftotal, _ := hex.DecodeString(p.hash)
	hmac, _ := blake2b.New256(buftotal)
	hmac1, _ := hex.DecodeString(p.hmac)

	hmac.Write(buftotal)

	return bytes.Compare(hmac.Sum(nil), hmac1) == 0

}

type Proof struct {
	hash       string
	hmac       string
	isVerified bool
}

type MediaProof struct {
	Proof
	MediaType string
	// we convert the image to greyscale so we have less information to deal with
	bwImgResized   *image.RGBA
	bwImgGenerated *image.RGBA
}

type MediaOrganism struct {
	DNA     *image.RGBA
	Fitness int64
}

func (m *MediaProof) LoadAndConvert(imagePath string) *image.RGBA {
	imgFile, err := os.Open(imagePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("error loading image")
	}
	output, _ := os.Create("./resized.png")
	defer output.Close()
	img, _ := png.Decode(imgFile)
	convert := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/30, img.Bounds().Max.Y/30))
	draw.NearestNeighbor.Scale(convert, convert.Rect, img, img.Bounds(), draw.Over, nil)
	if err != nil {
		panic(err.Error())
	}
	_ = png.Encode(output, convert)

	imgFile2, err := os.Open("./resized.png")
	defer imgFile2.Close()
	img2, _, err := image.Decode(imgFile2)
	m.bwImgResized = img2.(*image.RGBA)
	return img2.(*image.RGBA)
}

func GetMediaProof(imagePath string) {

	M := new(MediaProof)

	rand.Seed(time.Now().UTC().UnixNano())
	target := M.LoadAndConvert(imagePath)

	pixels := target.Pix[:]
	buf := new(bytes.Buffer)

	_, err := buf.Write([]byte{0x01, 0x11})
	if err != nil {
		panic(err)
	}
	hash := sha3.New256()

	hash.Write(pixels)
	buf.Write(hash.Sum(nil))
	M.hash = hex.EncodeToString(buf.Bytes())
	hm := hmac.New(sha256.New, buf.Bytes())
	M.hmac = hex.EncodeToString(hm.Sum(nil))
	M.bwImgResized = nil
	M.bwImgGenerated = nil
	spew.Dump(M)

}
