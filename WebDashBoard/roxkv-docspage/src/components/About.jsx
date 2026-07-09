import React from 'react';
import { about } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function About() {
  return (
    <SectionWrapper id="about">
      <div className="max-w-3xl mx-auto glass rounded-2xl p-8 sm:p-12 glow-border text-center">
        <h2 className="text-2xl sm:text-3xl font-bold gradient-text">{about.heading}</h2>
        <div className="w-16 h-1 rounded-full bg-gradient-to-r from-purple-500 to-purple-400 mx-auto mt-4" />
        <p className="text-[#9898ab] text-base sm:text-lg leading-relaxed mt-6">
          {about.description}
        </p>
      </div>
    </SectionWrapper>
  );
}
