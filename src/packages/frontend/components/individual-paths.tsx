import Image from "next/image";
import React from "react";

// Define a type for the card data
type CardData = {
  name: string;
  description: string;
  imageSrc: string;
  borderColor: string;
};

// Example array of card data (you would fetch or pass this in real use cases)
const cards: CardData[] = [
  {
    name: "Card 1",
    description: "This is the first card",
    imageSrc: "/sample.jpg",
    borderColor: "border-red-500",
  },
  {
    name: "Card 2",
    description: "This is the second card",
    imageSrc: "/sample.jpg",
    borderColor: "border-red-500",
  },
  {
    name: "Card 3",
    description: "This is the third card",
    imageSrc: "/sample.jpg",
    borderColor: "border-red-500",
  },
];

const Card: React.FC<{ data: CardData }> = ({ data }) => {
  return (
    <div
      className={`container flex flex-row items-center rounded overflow-hidden shadow-lg bg-white hover:bg-gray-200 hover:cursor-pointer hover:shadow-2xl transition-shadow duration-200 ${data.borderColor}`}
    >
      <div>
        {/* Adjust the height as needed */}
        <Image
          src={data.imageSrc}
          alt={`Picture of ${data.name}`}
          objectFit="cover"
          className="rounded-l"
          width={80}
          height={80}
        />
      </div>
      <div className="px-5">
        <div className="font-bold text-xl mb-2 text-black">{data.name}</div>
        <p className="text-gray-700 text-base">{data.description}</p>
      </div>
    </div>
  );
};

const CardGrid: React.FC = () => {
  return (
    <div className="p-5 bg-blue-200">
      <div className="grid grid-cols-1 gap-4">
        {cards.map((card, index) => (
          <Card key={index} data={card} />
        ))}
      </div>
    </div>
  );
};

export default CardGrid;
