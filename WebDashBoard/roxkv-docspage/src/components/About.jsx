import React from 'react';
import { about } from '../data/content';
import { ExternalLink, UserCircle2 } from 'lucide-react';
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
        <div className="flex flex-col sm:flex-row justify-center gap-4 mt-8">
          <a
            href={about.repoUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-primary"
          >
            <ExternalLink className="w-4 h-4" />
            View Project on GitHub
          </a>
          <a
            href={about.profileUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary"
          >
            <UserCircle2 className="w-4 h-4" />
            {about.author}
          </a>
        </div>
      </div>
    </SectionWrapper>
  );
}
