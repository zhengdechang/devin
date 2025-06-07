import Hero from '@/components/Hero';
import NubankSection from '@/components/NubankSection';
import CollaborateSection from '@/components/CollaborateSection';
import UseCasesSection from '@/components/UseCasesSection';
import LearnTogetherSection from '@/components/LearnTogetherSection';
import ToolsSection from '@/components/ToolsSection';
import IndustryLeadersSection from '@/components/IndustryLeadersSection';

export default function Home() {
  return (
    <main>
      <Hero />
      <NubankSection />
      <CollaborateSection />
      <UseCasesSection />
      <LearnTogetherSection />
      <ToolsSection />
      <IndustryLeadersSection />
      {/* You can add other sections/components below the Hero component if needed */}
    </main>
  );
}
