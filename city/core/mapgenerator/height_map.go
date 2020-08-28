package mapgenerator

import (
	"math"
	"math/rand"
)

var (
	DefaultGenerateConfig = HeightMapGenerateConfig{
		drops:        60,
		passes:       4,
		stability:    3,
		particlesMin: 400,
		particlesMax: 800,
	}
)

type HeightMap struct {
	Width  int
	Height int
	data   []int
	config HeightMapGenerateConfig
}

type HeightMapGenerateConfig struct {
	drops        int
	passes       int
	stability    int
	particlesMin int
	particlesMax int
}

func Generate(width, height int) *HeightMap {
	height_map := &HeightMap{config: DefaultGenerateConfig, data: make([]int, width*height), Width: width, Height: height}

	if height < 8 || width < 8 {
		// map too small do nothing
		return height_map
	}

	for i := 0; i < height_map.config.passes; i++ {
		drop_x, drop_y := width/2, height/2
		height_map.dropParticles(drop_x, drop_y, i)
		height_map.nextGenerateConfig()
	}

	GaussianFilter(height_map)

	return height_map
}

func (m *HeightMap) Value(pos_x, pos_y int) int {
	return m.data[pos_x+pos_y*m.Width]
}

func (m *HeightMap) dropParticles(pos_x, pos_y, interation int) {
	for i := 0; i < m.config.drops; i++ {
		particles := rand.Intn(m.config.particlesMax-m.config.particlesMin) + m.config.particlesMin

		for j := 0; j < particles; j++ {
			final_pos := m.checkFinalPos(pos_x, pos_y)
			m.data[final_pos] += 1
		}

		if interation == 0 {
			pos_x = rand.Intn(m.Width-8) + 4
			pos_y = rand.Intn(m.Height-8) + 4
		} else {
			pos_x = rand.Intn(m.Width/2) + m.Width/4
			pos_y = rand.Intn(m.Height/2) + m.Height/4
		}
	}
}

func (m *HeightMap) checkFinalPos(pos_x, pos_y int) int {
	pos := pos_x + pos_y*m.Width
	stability := m.config.stability

	if m.data[pos] == 0 {
		return pos
	}

	neighbors := []int{}
	for i := pos_x - stability; i < pos_x+stability; i++ {
		for j := pos_y - stability; j < pos_y+stability; j++ {
			neighbor_pos := i + j*m.Width
			if neighbor_pos >= 0 && neighbor_pos < m.Width*m.Height && m.data[neighbor_pos] < m.data[pos] {
				neighbors = append(neighbors, neighbor_pos)
			}
		}
	}

	if len(neighbors) == 0 {
		return pos
	}
	rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
	return neighbors[0]
}

func (m *HeightMap) nextGenerateConfig() {
	m.config.drops /= 2
	// m.config.passes -= 1
	m.config.particlesMin = int(math.Floor(float64(m.config.particlesMin) * 1.1))
	m.config.particlesMax = int(math.Floor(float64(m.config.particlesMax) * 1.1))
}
