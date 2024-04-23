"use client";
import { cn } from '@/libs/utils';
import { AutoCompleteData } from '@/types/autocomplete';
import Image from 'next/image';
import { useEffect, useRef, useState } from 'react';
import { Input, InputProps } from './ui/input';

interface AutoCompleteInputProps extends InputProps {
    data: AutoCompleteData[];
    onSelectValue: (url: string) => void;
}

const AutoCompleteInput: React.FC<AutoCompleteInputProps> = ({ data, onSelectValue, ...props }) => {
    const [selected, setSelected] = useState("");
    const [isOpen, setIsOpen] = useState(true);
    const dropdownRef = useRef<HTMLUListElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);

    // Handle input change
    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = event.target.value;
        setSelected(value);
        if (props.onChange) {
            props.onChange(event);
        }
    };

    // Handle item click
    const handleItemClick = (item: AutoCompleteData) => {
        onSelectValue(item.url);
        setSelected(item.title);
        if (inputRef && 'current' in inputRef && inputRef.current) {
            (inputRef.current as HTMLInputElement).focus();
        }
        setIsOpen(false);
    };


    // Handle click outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node) && inputRef.current && !inputRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };

        document.addEventListener('click', handleClickOutside);

        return () => {
            document.removeEventListener('click', handleClickOutside);
        };
    }, [])

    // Update open state when value changes
    useEffect(() => {
        if (props.value) {
            setIsOpen(true);
        }
    }, [props.value]);

    return (
        <div className="flex flex-col flex-1 relative">
            {/* Text Input */}
            <Input
                {...props}
                ref={inputRef}
                value={selected}
                onChange={handleInputChange}
                onFocus={() => setIsOpen(true)}
            />

            {/* Content */}
            {isOpen && data && data.length>0&& (
                <ul
                    ref={dropdownRef}
                    className="rounded-lg absolute top-16 bg-white text-black max-h-[500px] custom-scrollbar overflow-y-auto overflow-x-hidden w-full"
                    role="listbox"
                    aria-labelledby="autocomplete-input"
                >
                    {/* Mapping Data */}
                    {data.map((item) => (
                        <li
                            key={item.pageid}
                            role="option"
                            aria-selected={selected === item.title}
                            tabIndex={0}
                            onClick={() => handleItemClick(item)}
                            onKeyDown={(e) => e.key === 'Enter' && handleItemClick(item)}
                            className={cn(
                                "flex items-center gap-5 px-6 py-4 hover:bg-zinc-200 hover:cursor-pointer",
                                selected === item.title && "bg-slate-100"
                            )}
                        >
                            {/* Image display */}
                            <Image
                                src={item.image || '/default.png'}
                                width={80}
                                height={80}
                                alt={item.title}
                                sizes="80px"
                                className="aspect-square w-20 object-cover rounded-lg flex-shrink-0"
                            />
                            {/* Text Content */}
                            <div className="flex flex-col gap-2">
                                <h3 className="text-lg font-bold">{item.title}</h3>
                                <p className="text-sm line-clamp-3">{item.description}</p>
                            </div>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

AutoCompleteInput.displayName = 'AutoCompleteInput';

export default AutoCompleteInput;
