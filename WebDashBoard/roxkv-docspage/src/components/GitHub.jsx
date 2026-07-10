import React from 'react';
import { github } from '../data/content';
import { GitBranch, Star, GitCommit, ExternalLink, BookOpen, UserCircle2 } from 'lucide-react';
import SectionWrapper from './SectionWrapper';

export default function GitHub() {
  return (
    <SectionWrapper id="github">
      <div className="max-w-2xl mx-auto glass rounded-2xl p-8 sm:p-10 glow-border text-center">
        <GitBranch className="w-12 h-12 text-purple-400 mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-white font-mono mt-4">{github.repoName}</h2>
        <p className="text-[#9898ab] mt-3 leading-relaxed text-sm sm:text-base">
          {github.description}
        </p>
        <div className="flex justify-center gap-6 mt-6">
          <div className="flex items-center gap-2 text-sm text-[#9898ab]">
            <Star className="w-4 h-4" /> {github.stars}
          </div>
          <div className="flex items-center gap-2 text-sm text-[#9898ab]">
            <GitCommit className="w-4 h-4" /> {github.commits} commits
          </div>
        </div>
        <div className="flex flex-col sm:flex-row gap-4 justify-center mt-8">
          <a href={github.repoUrl} target="_blank" rel="noopener noreferrer" className="btn-primary">
            <ExternalLink className="w-4 h-4" /> View Repository
          </a>
          <a href={github.docsUrl} target="_blank" rel="noopener noreferrer" className="btn-secondary">
            <BookOpen className="w-4 h-4" /> Read Documentation
          </a>
        </div>
        <a
          href={github.profileUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 mt-6 text-sm text-[#9898ab] hover:text-purple-300 transition-colors"
        >
          <UserCircle2 className="w-4 h-4" />
          Built by @{github.ownerName}
        </a>
      </div>
    </SectionWrapper>
  );
}
