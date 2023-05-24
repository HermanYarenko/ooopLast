package main

import (
	"fmt"
	"math"
	"math/rand"

	//"sort"
	"strings"
	"time"
)

const SIDE_SIZE int = 10

var Names = []string{"Pudge", "Pugna", "Pangolier", "Aparet", "Alcgemist", "AntiMage", "Abadon", "Artes", "Omniknate", "Eminem",
	"Mia", "Meepo", "Scorpion", "ChaosKnite", "Lucan", "Amura", "Monkey", "Heart", "Duster",
	"Eve", "Alex", "Benzema", "Earth", "Strom", "Ember", "Void", "Titan", "Slick",
	"Dizel", "Murder", "Terry", "Coby", "Jordan", "Waen", "Deadpool", "Bmw", "Jonh",
	"Messi", "Mad", "Jordan", "Lebron", "Groot", "DarkSeer", "Zayac", "Dendi", "Pure",
	"Cent", "Tecgis", "Io", "Herman"}

type Battlefield struct {
	fields [][]string
}

func NewBattlefield(size int) *Battlefield {
	fields := make([][]string, size)
	for i := 0; i < size; i++ {
		fields[i] = make([]string, size)
		for y := 0; y < size; y++ {
			fields[i][y] = "-----"
		}
	}
	return &Battlefield{fields: fields}
}

func (bf *Battlefield) isEmpty(pos Point2D) bool {
	return bf.fields[pos.y][pos.x] == "-----"
}

func (bf *Battlefield) clearField(pos Point2D) {
	bf.fields[pos.y][pos.x] = "-----"
}

func (bf *Battlefield) setField(pos Point2D) {
	bf.fields[pos.y][pos.x] = "--X--"
}

type Point2D struct {
	x int
	y int
}

// Полный конструктор
func NewPoint2D(x int, y int) *Point2D {
	return &Point2D{x: x, y: y}
}

// Конструктор без параметров
func NewEmptyPoint2D() *Point2D {
	return &Point2D{x: 0, y: 0}
}

func (p *Point2D) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

func (p *Point2D) IsEqual(pos *Point2D) bool {
	if pos.y == p.y && pos.x == p.x {
		return true
	}
	return false
}

func (p *Point2D) GetDistance(other *Point2D) float64 {
	return math.Sqrt(math.Pow(float64(p.x-other.x), 2) + math.Pow(float64(p.y-other.y), 2))
}

type BaseHero struct {
	name           string
	herotype       string
	health         float64
	healthMax      int
	attack         int
	defense        int
	speed          int
	attackDistance float64
	damageMin      int
	damageMax      int
	posX           int
	posY           int
	side           string
	state          int
}

func NewBaseHero(name, herotype string, health float64, attack, defense, speed, damageMin, damageMax, posX, posY int) *BaseHero {
	h := &BaseHero{
		name:           name,
		herotype:       herotype,
		health:         health,
		healthMax:      int(health),
		attack:         attack,
		defense:        defense,
		speed:          speed,
		attackDistance: math.Sqrt(2),
		damageMin:      damageMin,
		damageMax:      damageMax,
		side:           "",
		state:          1,
		posX:           posX,
		posY:           posY,
	}
	return h
}

func NewHero(name string, x, y int) *BaseHero {
	return NewBaseHero(name, "Monk", 2, 3, 4, 1, 2, 1, x, y)
}

// Вывод всех полей в строковом виде
func (hero BaseHero) String() string {
	return fmt.Sprintf("%s:%s: health: %.1f/%d ♥, %s, attack: %d ⚔, defense: %d ✊, speed: %d, damage: %d-%d ☠, pos: %d,%d, state: %d",
		hero.side,
		hero.name,
		hero.health,
		hero.healthMax,
		hero.herotype,
		hero.attack,
		hero.defense,
		hero.speed,
		hero.damageMin,
		hero.damageMax,
		hero.posX,
		hero.posY,
		hero.state)
}

func (hero BaseHero) getPos() *Point2D {
	return &Point2D{hero.posX, hero.posY}
}

// Сортировка по убыванию скорости
func (hero BaseHero) compareTo(otherHero BaseHero) int {
	if hero.speed < otherHero.speed {
		return 1
	} else if hero.speed > otherHero.speed {
		return -1
	}
	return 0
}

// Рассчет - Расстояние - от текущей позиции до цели
func (hero BaseHero) getDistance(otherHero *BaseHero) float64 {
	return math.Sqrt(math.Pow(float64(otherHero.posX-hero.posX), 2) + math.Pow(float64(otherHero.posY-hero.posY), 2))
}

// Рассчет - Урон - по выбранной цели, минимальный урон == 1
func (hero BaseHero) getDamage(otherHero *BaseHero) float64 {
	damage := float64(otherHero.defense - hero.attack - ((hero.damageMin + hero.damageMax) / 2))
	if damage < 0 {
		return damage * -1
	}
	return 1
}

// Рассчет - Ближайшая живая цель - если все умерли, то последняя мертвая цель
func (hero BaseHero) getTarget(enemySide []*BaseHero) *BaseHero {
	target := enemySide[len(enemySide)-1]
	minDistance := math.Abs(hero.getDistance(target))
	for i := 0; i < len(enemySide)-1; i++ {
		enemy := enemySide[i]
		if enemy.state != -1 {
			newDistance := math.Abs(hero.getDistance(enemy))
			if newDistance < minDistance {
				minDistance = newDistance
				target = enemy
			} else if enemySide[len(enemySide)-1].state == -1 {
				minDistance = newDistance
				target = enemy
			}
		}
	}
	return target
}

// Эффект - Получение урона
func (hero *BaseHero) doDamage(damage float64) {
	hero.health -= damage
	if hero.health <= 0 {
		hero.health = 0
		hero.state = -1
		fmt.Println(hero)
	} else {
		fmt.Println(hero)
	}
}

func (attacker *BaseHero) doAttack(battlefield *Battlefield, defender *BaseHero) {
	damage := attacker.getDamage(defender)
	fmt.Printf("%s получает урон %d от %s:%s\n", defender, damage, attacker.side, attacker.name)
	defender.doDamage(damage)
	if defender.state == -1 {
		battlefield.clearField(*defender.getPos())
	}
}

func (mover *BaseHero) doMoveTo(battlefield *Battlefield, target *BaseHero) {
	distanceX := target.posX - mover.posX
	distanceY := target.posY - mover.posY
	if distanceX == 0 {
		distanceX = 1
	}
	if distanceY == 0 {
		distanceY = 1
	}
	x := distanceX / int(math.Abs(float64(distanceX)))
	y := distanceY / int(math.Abs(float64(distanceY)))

	if math.Abs(float64(distanceX)) > 1 && math.Abs(float64(distanceY)) > 1 &&
		battlefield.isEmpty(Point2D{mover.posX + x, mover.posY + y}) {
		battlefield.setField(Point2D{mover.posX + x, mover.posY + y})
		battlefield.clearField(*mover.getPos())
		mover.posY = mover.posX + x
		mover.posY = mover.posY + y
	} else if math.Abs(float64(distanceX)) > 1 &&
		battlefield.isEmpty(Point2D{mover.posX + x, mover.posY}) {
		battlefield.setField(Point2D{mover.posX + x, mover.posY})
		battlefield.clearField(*mover.getPos())
		mover.posY = mover.posX + x
	} else if battlefield.isEmpty(Point2D{mover.posX, mover.posY + y}) {
		battlefield.setField(Point2D{mover.posX, mover.posY + y})
		battlefield.clearField(*mover.getPos())
		mover.posY = mover.posY + y
	}
	//fmt.Printf("%s -<- Переместился, на позицию %s\n", mover, mover.posX, mover.posY)
}

func (hero *BaseHero) doStep(battlefield *Battlefield, enemySide []*BaseHero) {
	if hero.state == 1 {
		target := hero.getTarget(enemySide)
		if target.state != -1 {
			if math.Abs(float64(hero.getDistance(target))) <= hero.attackDistance {
				hero.doAttack(battlefield, target)
				return
			} else {
				hero.doMoveTo(battlefield, target)
				return
			}
		}
		fmt.Printf("У %s -<- Нет цели для атаки, на поле нет живых противников\n", hero)
	}
}

func main() {
	RadiantSide := make([]*BaseHero, SIDE_SIZE)
	DireSide := make([]*BaseHero, SIDE_SIZE)
	allSide := make([]*BaseHero, 0)

	battlefield := NewBattlefield(SIDE_SIZE)

	rand.Seed(time.Now().UnixNano())

	x := 0
	for i := 0; i < SIDE_SIZE; i++ {
		switch rand.Intn(4) {
		case 0:
			RadiantSide[i] = NewHero(getName(), x, i)
		case 1:
			RadiantSide[i] = NewHero(getName(), x, i)
		case 2:
			RadiantSide[i] = NewHero(getName(), x, i)
		default:
			RadiantSide[i] = NewHero(getName(), x, i)
		}
		battlefield.fields[RadiantSide[i].posY][RadiantSide[i].posX] = "--X--"
		RadiantSide[i].side = "RadiantSide"
	}

	x = SIDE_SIZE - 1
	for i := 0; i < SIDE_SIZE; i++ {
		switch rand.Intn(4) {
		case 0:
			DireSide[i] = NewHero(getName(), x, i)
		case 1:
			DireSide[i] = NewHero(getName(), x, i)
		case 2:
			DireSide[i] = NewHero(getName(), x, i)
		default:
			DireSide[i] = NewHero(getName(), x, i)
		}
		battlefield.fields[DireSide[i].posY][DireSide[i].posX] = "--X--"
		DireSide[i].side = "DireSide"
	}

	allSide = append(allSide, RadiantSide...)
	allSide = append(allSide, DireSide...)

	input := ""

	fmt.Println(" PowerShell не хочет краситься ")
	for {
		if strings.Compare(input, "ecs") == 0 {
			break
		}
		for _, hero := range allSide {
			if strings.Compare(hero.side, "whiteSide") == 0 {
				hero.doStep(battlefield, DireSide)
			} else {
				hero.doStep(battlefield, RadiantSide)
			}
		}
		fmt.Println("\n" + strings.ReplaceAll(fmt.Sprintf("%v", battlefield.fields), "], ", "]\n"))
		fmt.Println("\nвведите 'ENTER' для след.шага() или 'ecs' чтоб завершить игру ")
		fmt.Scan(&input)
	}
}

// Выбор случайного имени для героя
func getName() string {
	return Names[rand.Intn(len(Names))]
}
