import Navbar from '../components/Navbar';
import Hero from '../components/Hero';
import Stats from '../components/Stats';
import WhatIsRoxAI from '../components/WhatIsRoxAI';
import WhyItExists from '../components/WhyItExists';
import Capabilities from '../components/Capabilities';
import DashboardShowcase from '../components/DashboardShowcase';
import Architecture from '../components/Architecture';
import Engineering from '../components/Engineering';
import TechStack from '../components/TechStack';
import Demo from '../components/Demo';
import GitHub from '../components/GitHub';
import RunLocally from '../components/RunLocally';
import Gallery from '../components/Gallery';
import Roadmap from '../components/Roadmap';
import About from '../components/About';
import Footer from '../components/Footer';

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-[#0a0a0f] text-[#f1f1f4]">
      <Navbar />
      <Hero />

      <div className="section-divider" />
      <Stats />

      <div className="section-divider" />
      <WhatIsRoxAI />

      <div className="section-divider" />
      <WhyItExists />

      <div className="section-divider" />
      <Capabilities />

      <div className="section-divider" />
      <DashboardShowcase />

      <div className="section-divider" />
      <Architecture />

      <div className="section-divider" />
      <Engineering />

      <div className="section-divider" />
      <TechStack />

      <div className="section-divider" />
      <Demo />

      <div className="section-divider" />
      <GitHub />

      <div className="section-divider" />
      <RunLocally />

      <div className="section-divider" />
      <Gallery />

      <div className="section-divider" />
      <Roadmap />

      <div className="section-divider" />
      <About />

      <Footer />
    </div>
  );
}
