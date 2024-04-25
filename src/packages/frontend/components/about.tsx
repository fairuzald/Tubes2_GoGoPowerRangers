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
import * as React from "react";
import { useState, useEffect } from "react";

export function DrawerAbout() {
  const [goal, setGoal] = React.useState(350);

  function onClick(adjustment: number) {
    setGoal(Math.max(200, Math.min(400, goal + adjustment)));
  }

  const [hasScrolled, setHasScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => {
      // Set the state to true if scrolled down by more than 50px, for instance
      setHasScrolled(window.scrollY > 10);
    };

    // Attach the listener to the window scroll event
    window.addEventListener("scroll", onScroll);

    // Clean up the listener when the component unmounts
    return () => {
      window.removeEventListener("scroll", onScroll);
    };
  }, []);

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button
          variant="link"
          className={`text-xl cursor-pointer hover:font-bold shadow-none hover:shadow-lg p-0 ${
            hasScrolled ? "hover:text-black " : "hover:text-yellow-hover"
          }`}
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
            <div className="flex flex-col items-center lg:flex-row my-4 gap-4">
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
                  Anda dapat menggunakan website WikiRace dengan cara: <br /> 1.
                  Masukkan judul pada source dan destination
                  <br /> 2. Pilih metode pencarian dengan algoritma Iterative
                  Deepening Search (IDS) atau Breadth First Search (BFS)
                  <br /> 3. Tekan tombol {"Go!"} untuk memulai pencarian <br />
                  4. Hasil berupa tampilan box dan visualisasi berupa graf
                </DrawerDescription>
              </div>
            </div>
          </DrawerHeader>
          <DrawerFooter>
            <DrawerClose asChild>
              <Button
                variant="outline"
                className="text-lg hover:bg-white hover:text-yellow-primary"
              >
                Back
              </Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </DrawerContent>
    </Drawer>
  );
}
