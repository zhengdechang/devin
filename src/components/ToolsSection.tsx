import React from 'react';

interface ToolLogoProps {
  name: string;
}

const ToolLogo: React.FC<ToolLogoProps> = ({ name }) => {
  return (
    <div className="bg-gray-200 p-4 h-20 rounded-lg flex items-center justify-center shadow">
      <p className="text-gray-700 font-medium text-sm">{name}</p>
    </div>
  );
};

const ToolsSection = () => {
  const toolNames = [
    "Confluence", "Airtable", "Segment", "Asana", "Notion", "Stripe",
    "AWS", "GitHub", "Datadog", "Linear", "Databricks", "Slack",
    "Google Drive", "Sentry", "PostgreSQL", "Azure", "Snowflake", "MongoDB"
    // Add more if needed, or indicate that many more are supported
  ];

  return (
    <section className="py-12 md:py-24 bg-white">
      <div className="container mx-auto px-6">
        <h2 className="text-3xl md:text-4xl font-bold text-center text-gray-800 mb-4">
          Able to work with hundreds of tools
        </h2>
        <p className="text-lg text-gray-600 text-center max-w-2xl mx-auto mb-16">
          Devin seamlessly integrates with your existing stack. Here are just a few examples of tools it can work with:
        </p>
        <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-4 md:gap-6">
          {toolNames.map((toolName) => (
            <ToolLogo key={toolName} name={toolName} />
          ))}
        </div>
        <p className="text-center text-gray-500 mt-8">
          ...and many more, adapting to your project's specific needs.
        </p>
      </div>
    </section>
  );
};

export default ToolsSection;
