import * as React from "react";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import InteractiveImage from "@/components/ui/interactive-image";

export function DrawerAbout() {
  const [goal, setGoal] = React.useState(350);

  function onClick(adjustment: number) {
    setGoal(Math.max(200, Math.min(400, goal + adjustment)));
  }

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button
          variant="link"
          className={`text-xl cursor-pointer hover:font-bold shadow-none hover:shadow-lg p-0`}
        >
          About
        </Button>
      </DrawerTrigger>
      <DrawerContent>
        <div className="mx-auto w-full">
          <DrawerHeader>
            <DrawerTitle>
              Pemanfaatan Algoritma IDS dan BFS dalam Permainan WikiRace
            </DrawerTitle>
            <div className="flex flex-row my-4 gap-4">
              <InteractiveImage />
              <div className="flex flex-col gap-4">
                <DrawerDescription>
                  WikiRace atau Wiki Game adalah permainan yang melibatkan
                  Wikipedia, sebuah ensiklopedia daring gratis yang dikelola
                  oleh berbagai relawan di dunia, dimana pemain mulai pada suatu
                  artikel Wikipedia dan harus menelusuri artikel-artikel lain
                  pada Wikipedia (dengan mengeklik tautan di dalam setiap
                  artikel) untuk menuju suatu artikel lain yang telah ditentukan
                  sebelumnya dalam waktu paling singkat atau klik (artikel)
                  paling sedikit.{" "}
                </DrawerDescription>
                <DrawerDescription>
                  <ol>
                    <li>
                      Anda dapat menggunakan website WikiRace dengan cara:
                    </li>
                    <li>
                      1. Masukkan judul pada input {"from"} (source) dan input{" "}
                      {"to"} (destination)
                    </li>
                    <li>
                      2. Pilih metode pencarian dengan algoritma Iterative
                      Deepening Search (IDS) atau Breadth First Search (BFS){" "}
                    </li>
                    <li>3. Tekan tombol {"Go!"} untuk memulai pencarian</li>
                    <li>
                      4. Hasil dapat berupa tampilan box yang dapat Anda klik
                      dan visualisasi berupa graf
                    </li>
                  </ol>
                </DrawerDescription>
              </div>
            </div>
          </DrawerHeader>
          <DrawerFooter>
            <DrawerClose asChild>
              <Button variant="outline" className="text-lg">
                Back
              </Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </DrawerContent>
    </Drawer>
  );
}
