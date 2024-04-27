# WikiRace Application

<p align="center">
  <a href="https://tubes2-go-go-power-rangers.vercel.app/">
    <h3 align="center">WikiRace Application</h3>
  </a>
</p>

<p align="center">Next.js boilerplate that uses <a href="https://pkg.go.dev/github.com/gin-gonic/gin/">GIN</a> as the API backend.</p>

<br/>



## Table of Contents

- [General Information](#general-information)
- [Demo](#demo)
- [Features](#features)
- [Getting Started](#getting-started)
- [Learn More](#learn-more)
- [Deploy on Vercel](#deploy-on-vercel)
- [Team Members](#team-members)
- [Program Structure](#program-structure)

## General Information

WikiRace or Wiki Game is a game that involves Wikipedia, a free online encyclopedia managed by various volunteers around the world, where players start at a specific Wikipedia article and must navigate through other Wikipedia articles (by clicking on links within each article) to reach another predetermined article in the shortest time or with the fewest clicks (articles).

![Example screenshot](https://miro.medium.com/v2/resize:fit:640/format:webp/1*NwVK37pwD5dHFNBfevfo1w.png)

## Demo

Check out the live demo: [Wiki Race Demo](https://tubes2-go-go-power-rangers.vercel.app)

## Features

| **No.** |                                               **Poin**                                                | **Ya** | **Tidak** |
| :-----: | :---------------------------------------------------------------------------------------------------: | :----: | --------- |
|   1.    |            Program menyediakan pencarian dengan algoritma Iterative Deepening Search (IDS)            |   ✓    |           |
|   2.    |               Program menyediakan pencarian dengan algoritma Breadth First Search (BFS)               |   ✓    |           |
|   3.    | (Bonus) Program dapat menampilkan seluruh rute terpendek (tidak hanya satu) dan memvisualisasikannya. |   ✓    |           |
|   4.    |           (Bonus) Program dijalankan menggunakan Docker baik untuk frontend maupun backend.           |   ✓    |           |
|   5.    |                     (Bonus) Rute terpendek dengan durasi kurang dari satu menit.                      |   ✓    |           |
|   6.    |                                        (Bonus) Video kelompok.                                        |   ✓    |           |


## Development

To run the application locally:

### Prerequisites

Before starting the development process, make sure you have the following software installed on your machine:

- [Node.js](https://nodejs.org/) (LTS version recommended)
- [Golang](https://go.dev/) (LTS version recommended)
- [npm](https://www.npmjs.com/) or [yarn](https://yarnpkg.com/) (for managing JavaScript dependencies)

### 1. Clone the Repository

```bash
git clone https://github.com/fairuzald/Tubes2_GoGoPowerRangers
cd Tubes2_GoGoPowerRangers
```

### 2. Navigate to the Source (src) Directory

```bash
cd src

```

### 3. Install All Frontend and Backend Dependencies

```bash
npm run install-all
# or
yarn install-all
# or
pnpm install-all
```

### 4. Run Next.js Development Server

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
```

### 5: Open your browser and navigate to

Client-side is running on [localhost:3000](http://localhost:3000), and the server is on [localhost:8080](http://localhost:8080).

To run the application using docker:


### 1. Clone the Repository

```bash
git clone https://github.com/fairuzald/Tubes2_GoGoPowerRangers
cd Tubes2_GoGoPowerRangers
```

### 2. Navigate to the Source (src) Directory

```bash
cd src

```

### 3. Build and run with docker

```bash
npm run start
# or
yarn start
# or
pnpm start
```
## Project Status

Project is complete


## Team Members

| **NIM**  |         **Nama**          |
| :------: | :-----------------------: |
| 13522042 |       Amalia Putri        |
| 13522057 | Moh Fairuz Alauddin Yahya |
| 13522080 |   Julian Chandra Sutadi   |

## Program Structure

```
.
└── Tubes2_GoGoPowerRangers
   ├── doc
   │ └── Tucil2_13522010_13522042.pdf
   |
   ├── src
   │ ├── node_modules
   │ └── packages
   |   ├── backend
   |   | └── ...
   │   └── frontend
   |     └── ...
   |
   └── README.md
```
