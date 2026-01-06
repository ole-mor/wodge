Received 11 July 2025, accepted 18 August 2025, date of publication 21 August 2025, date of current version 27 August 2025.
Digital Object Identifier 10.1109/ACCESS.2025.3601206
A Systematic Literature Review of Hallucinations
in Large Language Models
CHRISTIAN WOESLE , LEOPOLD FISCHER-BRANDIES , (Member, IEEE),
AND RICARDO BUETTNER , (Senior Member, IEEE)
Chair of Hybrid Intelligence, Helmut-Schmidt-University/University of the Federal Armed Forces Hamburg, 22043 Hamburg, Germany
Corresponding author: Ricardo Buettner (buettner@hsu-hh.de)
This work was supported by the Open-Access-Publication-Fund of the Helmut-Schmidt-University/University of the Federal Armed
Forces Hamburg.
ABSTRACT This review systematically maps research on hallucinations in large language models using
a descriptive scheme that links model outputs to four system architectures: unaugmented generation, post-
hoc reactive validation, proactive detection-and-mitigation, and fully integrated detection-and-mitigation
designs. Our methodology for this systematic review follows the PRISMA guidelines to ensure transparency
and reproducibility. We searched IEEE Xplore, ACM Digital Library, and ScienceDirect for studies
published between 2015 and January 2025 and extracted 125 peer-reviewed papers across nine application
domains. Quantitative analysis shows that question answering and multimodal tasks account for 48% of all
papers, whereas software engineering, educational technology, and autonomous systems are underexplored.
Although 87.5% of the studies rely on additional reactive or proactive defenses, only 8.8% implement
integrated architecture-level safeguards, revealing a critical gap in unified and dynamic architectures.
The resulting classification matrix and domain map provide a diagnostic tool for locating blind spots
and comparing architectural maturity. Three actionable priorities emerge: develop integrated reasoning-
and-verification loops that pre-empt hallucinations; transfer proven causal-intervention and multi-agent
validation pipelines to high-stakes, under-represented domains and benchmark them under real conditions;
and build modular, cross-domain evaluation frameworks that isolate the contribution of individual
mitigation components and support ablation studies. By consolidating fragmented evidence and quantifying
architecture-domain imbalances, this review establishes a traceable foundation for engineering reliable,
explainable, and domain-adaptable countermeasures to hallucinations in generative language technology.
INDEX TERMS Large language models, hallucinations, architecture, detection techniques, mitigation
strategies, systematic literature review.
I. INTRODUCTION
Large Language Models (LLMs) such as GPT-3, GPT-4, and
PaLM have demonstrated remarkable capabilities in natural
language understanding and generation [1], [2], question
answering [3], code generation [4], powering advancements
in dialogue systems [5], [6], and summarization [7], [8].
These models achieve their performance through large-scale
pretraining on massive textual corpora, learning statistical
patterns that allow them to produce fluent, contextually
The associate editor coordinating the review of this manuscript and
approving it for publication was Prakasam Periasamy .
VOLUME 13, 2025
appropriate language [9]. However, a growing body of
evidence has raised concerns about a fundamental limitation:
hallucinations.
Hallucinations are instances where LLMs generate content
that is grammatically fluent and syntactically coherent,
yet factually incorrect, unverifiable, or inconsistent with
a source input [10]. This issue manifests across nearly
all natural language generation tasks. In summarization,
models may fabricate facts not present in the source
document [11]; in question answering [9], they may con-
fidently present incorrect responses [9]; and in scientific
or medical domains, they may generate plausible-sounding
2025 The Authors. This work is licensed under a Creative Commons Attribution 4.0 License.
For more information, see https://creativecommons.org/licenses/by/4.0/ 148231
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
yet fabricated references [2]. Hsu and Roberts [2] found
that LLMs used in clinical NLP pipelines can introduce
annotations that appear consistent but lack verification,
raising concerns about factual integrity. Bzdok et al. [1] sim-
ilarly observed that LLM outputs often reflect ‘‘syntactically
elegant guesses,’’ particularly when models are trained on
unstructured or domain-general corpora.
Despite recent innovations, including retrieval-augmented
generation [12], Chain-of-Thought prompting [13], role-
playing multi-agent prompting [14], and semantic entropy-
based detection [10], the hallucination problem remains
unresolved. Current models prioritize fluency and likelihood
maximization over factual accuracy, and existing evaluation
metrics often fail to detect hallucinated content [9], [11].
Moreover, Ji et al. [9] emphasize that mitigation techniques
vary greatly in effectiveness across domains and models, and
that there is limited consensus in the literature regarding the
most effective or generalizable solutions.
Although the phenomenon of hallucination also occurs
in humans, the phenomenon of hallucination in LLMs has
not yet been systematically investigated from a unified
conceptual and methodological lens. Most existing works
either address hallucination as a side issue or focus narrowly
on task-specific metrics. While several recent works survey
hallucination phenomena in LLMs [15], [16], [17], covering
aspects such as causes [18], definitions [19], taxonomies [20]
and mitigation strategies [21] to varying extents, these efforts
remain limited in scope or depth, and none attempt a unified
synthesis across these dimensions. For example, Huang et al.
[20] provide a comprehensive taxonomy but do not ground
their approach in a formal theoretical framework. Reddy et al.
[15] and Perković et al. [16] offer useful overviews, yet
do not adopt a unifying conceptual framework. Maleki et al.
[19], on the other hand, focus on a meta-level analysis of
how the term ‘hallucination’ is used across different AI
subfields.
This paper addresses this gap through a systematic
review of the literature on hallucinations in LLMs, ana-
lyzing 125 peer-reviewed publications within a descriptive
organizational scheme based on system architectures. The
primary objective is to categorize and evaluate existing
strategies for detection and mitigation, identifying their
applicability and architectural sophistication across vari-
ous domains such as Question Answering & Informa-
tion Retrieval, Multimodal & Vision-Language Models,
and Biomedical & Scientific Applications. A secondary
objective is to examine how this architectural framework
provides a useful lens for understanding hallucination
behaviors in artificial systems and for future design
principles.
We aim to answer the following research questions:
• How can system architecture categories be used to
analyze and categorize hallucinations in LLMs?
• In what ways can this architectural framework support
the classification of existing hallucination mitigation
strategies?
FIGURE 1. 125 papers across 9 domains were analyzed and mapped to
system architecture categories;only 8.8% reached Integrated LLM-DMS
Architectures.
• Which hallucination detection and mitigation approaches
demonstrate distinct architectural sophistication, and
what characterizes them as particularly promising?
• Which domains are currently underrepresented in hallu-
cination research, and how does their coverage compare
in terms of architectural sophistication?
To address these questions, we performed a domain-specific
classification of the hallucination literature and proposed a
multiaxial evaluation matrix. We considered papers published
between 2015, and January 17, 2025, reporting empirical
results, methodological innovations, or theoretical insights
related to hallucinations in LLMs. Our methodology follows
the PRISMA [22] guidelines to ensure transparency and
reproducibility.
By consolidating scattered results and linking them
to architectural frameworks, this work provides a solid
foundation for future research on robust, explainable, and
architecturally sound language models.
To address the identified research gap, this article presents
a systematic literature review (SLR) on hallucinations in
LLMs, combining empirical categorization with an architec-
turally grounded analytical lens (Figure 1). Drawing from our
proposed system architecture categories, 125 peer-reviewed
studies are evaluated along two key axes: their domain
of application and the sophistication of their hallucination
detection and mitigation strategies. This structured synthesis
provides a comprehensive overview of the current research
landscape, identifying dominant areas of focus, emerging
theoretical advancements, and significant gaps. The proposed
classification matrix serves both as a conceptual framework
and a diagnostic tool to support the development of more
explainable, architecturally aligned, and domain-adaptable
solutions to hallucination phenomena in generative AI
systems.
With regard to the research gap, the main contributions of
this work are as follows:
• Introduction of a novel descriptive organizational
scheme based on system architectures, enabling a struc-
tured analysis of hallucination detection and mitigation
strategies in LLMs across diverse domains.
• Systematic mapping of 125 peer-reviewed publications,
revealing a concentration of research in domains such
as Question Answering & Information Retrieval and
148232 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
Multimodal & Visual-Language Models, while identi-
fying substantial gaps in Software Engineering & Code
Generation, Educational AI & Learning Systems, and
Autonomous Systems & Robotics.
• Evaluation of architectural sophistication across
reviewed techniques, highlighting that only 8.8%
of studies exhibit Integrated LLM-Detection and
Mitigation System (DMS) Architectures, indicating the
need for more dynamic and unified system designs in
future research.
This work is structured as follows: In section II, prelimi-
nary concepts and fundamentals are elaborated to give a brief
introduction to the topic. Next, section III, describes our SLR
methodology. The results from our SLR are then presented in
section IV and consequently discussed in section V. Finally,
we conclude in section VI by highlighting the theoretical and
practical implications of our findings.
related topics but are not considered systematic reviews.
These works do not follow a standardized review methodol-
ogy, such as PRISMA or domain-specific taxonomies, and
typically do not provide transparent inclusion criteria or
structured synthesis of reviewed studies. However, they offer
valuable information about the broader discourse around
hallucinations and limitations in LLMs. Zhang et al. [21]
provide a comprehensive study of safety, security, and privacy
issues in LLMs, including hallucinations, and propose a
refined classification framework to distinguish these con-
cerns and their corresponding defenses. Sanu et al. [18]
present a general critique of LLM vulnerabilities, including
hallucinations, bias, adversarial attacks, and outdatedness,
framed as ‘loopholes’ in their deployment. Reddy et al. [15]
offer a conceptual discussion on the nature and types of
hallucinations in LLMs, supported by real-world examples
but lacking empirical aggregation of research. Perković et al.
[16] focus on the technical and social dimensions of
hallucinations, such as data provenance and trust, while
proposing mitigation strategies from a systems engineering
perspective. Lastly, Huang et al. [20] discuss hallucinations
within a broader overview of generative AI’s limitations,
offering definitions and conceptual categories rather than
systematically analyzing existing work.
Although not meeting the formal review criteria, these arti-
cles contribute important perspectives and help contextualize
the growing interest in LLM hallucination research.
II. PRELIMINARY CONCEPTS AND FUNDAMENTALS
Hallucination in LLMs refers to the generation of content that
is fluent, syntactically correct, and semantically plausible,
but factually incorrect, unfounded, or completely fabricated.
They differ from conventional natural language generation
errors by a high degree of linguistic coherence that can
mask underlying factual inaccuracies [9], [10]. Ji et al. [9]
define hallucination as inaccurate generation and propose
a multidimensional taxonomy based on origin (intrinsic vs.
extrinsic), nature (factual, logical, linguistic), and manifesta-
tion (addition, omission, distortion).
Evaluating hallucinations remains difficult. Classical met-
rics like ROUGE or BERTScore fail to capture the factual
accuracy [11]. More recent efforts include entailment-based
metrics (e.g., FactCC), QA-based approaches (e.g., QAGS),
and semantic methods such as SummaC and semantic
entropy [10], [23], [24]. Yet these techniques vary in their
reliability across domains, and human evaluation remains the
gold standard despite cost and subjectivity. Hallucination also
manifests differently across tasks: abstractive summarization
systems may generate unsupported facts; dialogue models
can produce fictional information when context is missing;
and data-to-text systems frequently hallucinate content not
present in the structured input [9], [11].
Mitigation efforts aim to reduce hallucinations through
architectural or inference-level innovations. Retrieval-
augmented generation (RAG) incorporates external factual
knowledge to support grounded outputs [12], while Chain-
of-Thought prompting encourages intermediate reasoning
steps [13]. Despite these advances, no approach has fully
eliminated hallucinations in open-ended text generation
tasks.
In addition to fully systematic reviews, the broader
literature also contains relevant but less formally structured
analyses. While Table 1 presents structured and methodolog-
ically grounded literature reviews focused specifically on
hallucinations in LLMs, several additional articles address
A. EVOLUTION OF HALLUCINATIONS IN LARGE
LANGUAGE MODELS
The problem of hallucination has evolved in parallel with the
growing capabilities of LLMs. Early neural language models
such as GPT-2 and BART showed signs of hallucination,
particularly in abstract summarization [9], [11]. Originally,
these problems were attributed to overfitting or weak
attentional mechanisms [9], [11]. However, with the advent of
transformer-based models such as GPT-3, GPT-4 and PaLM,
hallucinations became clearer and better understood as a
systemic result of predicting the next token based on noisy,
incomplete, or unstructured training data [9], [25].
LLMs have also been shown to hallucinate in specialized
domains such as clinical NLP and neuroscience. Hsu
and Roberts [2] demonstrated that LLMs used in weak
supervision pipelines for clinical named entity recogni-
tion can introduce semantically plausible but unverified
labels, raising concerns about label reliability. Similarly,
Bzdok et al. [1] highlighted that LLMs may produce
plausible but ungrounded outputs, especially when trained
on unstructured or domain-unspecific data, contributing to
the broader understanding of hallucinations as coherent
yet misaligned model behaviors. Farquhar et al. [10]
proposed the semantic entropy framework to detect such
cases and showed that high variance in semantic intent
across resampled outputs correlates with hallucination
risk.
VOLUME 13, 2025 148233
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
TABLE 1. Overview of recent literature reviews regarding hallucinations in LLMs. The table displays each review’s authors, publication year, and the
number of papers included in the review. Additionally, a summary is given, and the applied review structure is provided.
As awareness of the issue grew, the research commu-
nity shifted from post-hoc detection to active mitigation.
Retrieval-augmented models like RAG introduced dynamic
grounding by retrieving relevant factual data during infer-
ence [12]. Chain-of-Thought prompting showed promise in
mathematical and commonsense reasoning tasks by mim-
icking deliberative thought processes [13]. Reinforcement
learning from human feedback (RLHF) also became a
prominent method for aligning outputs with user expectations
and reducing hallucination frequency [26].
Recent work has further expanded the scope of hallucina-
tion research into new domains. Zhang et al. [25] proposed
an Entity-Relation–based reinforcement learning framework
to reduce spatial hallucination in path planning, highlighting
that the issue extends beyond text to reasoning in struc-
tured environments. Despite these advances, hallucination
remains unsolved, not only in open-domain generation, where
prompts are vague and training data may be contradictory,
but also in structured tasks, as recent research has shown.
The evolution of hallucination research reflects an ongoing
tension between scale, fluency, and factual grounding.
systematic literature review on hallucinations in LLMs, and
(ii) a categorization framework for synthesizing the identified
papers. The former describes how relevant studies were
identified, screened, and included, while the latter discusses
the approach for classifying the selected studies in terms of
domain focus and architecture categories.
B. SYSTEMATIC LITERATURE REVIEW PROCEDURE
We conducted a SLR following the preferred reporting
items for systematic reviews and meta-analyses (PRISMA)
guidelines by Page et al. [22]. An SLR systematically
gathers, evaluates, and synthesizes relevant evidence on
a defined topic to establish the state of the art, identify
research gaps, and provide a solid foundation for further
investigation. In our case, the SLR was designed to address
the research questions outlined in the Introduction, with
a focus on hallucination phenomena in LLMs, evaluated
through a domain-based and an organisational-technical
categorization. To ensure a structured and transparent review,
we followed eight methodological steps adapted from the
PRISMA guidelines [22]. These steps are summarized in
Figure 2 and described in more detail below.
III. METHODOLOGY
A. OVERVIEW
This section outlines the methodological steps adopted for
this study, which comprises two primary components: (i) a
1) LITERATURE SEARCH STRATEGY
Following the presented framework, the next step of this SLR
was to define the search string and the criteria to exclude
148234 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
• The publication is considered ‘‘off-topic’’ by at least
two authors. This criterion ensures the relevance of the
included literature.
• The publication is written in a language other than
English.
In the fourth step, the databases used for the SLR
were selected and searched. The present study used three
databases: IEEE Xplore Digital Library, ACM Digital
Library, and ScienceDirect database. The search term com-
bination remained the same for all databases and was applied
to search titles, abstracts, and keywords. IEEE Xplore was
last accessed on 01/17/2025, ACM on 01/17/2025, and
ScienceDirect on 01/17/2025.
FIGURE 2. Overview of the process for obtaining the literature base as
derived from the PRISMA guidelines [22].
individual search results. Our search string consists of the
following components:
Title, Abstract & Keywords: (‘‘hallucination*’’) AND
(‘‘LLM*’’ OR ‘‘large language model*’’)
As the research goal of this article is to investigate
hallucinations in LLMs, abbreviated LLMs, the search term
was formulated to comprise three interrelated components:
hallucinations, large language models, and their variations.
To achieve this, the terms ‘‘hallucination*’’, ‘‘LLM*’’,
and ‘‘large language model*’’ were included in the search
string. Boolean operators were used to ensure accuracy
and completeness. The ‘‘AND’’ operator was used to
find only those papers that address all key dimensions
simultaneously, while the ‘‘OR’’ operator expanded the scope
by including synonymous terms and variations. In addition,
the use of the wildcard (*) enabled the inclusion of
plural forms and variations of the terms, further improving
coverage.
To ensure the quality of the identified work, the following
exclusion criteria were defined:
Exclusion Criteria
• The publication was published before 2015.
• The publication type is not a peer-reviewed conference
paper or journal article. To ensure academic rigor,
formats such as poster sessions, editorials, interviews,
commentaries, and research-in-progress papers are
excluded.
2) STUDY SELECTION AND EXTRACTION
Before screening, four articles from IEEE Xplore, one
article from ACM, and four articles from ScienceDirect were
excluded due to incorrect article types. Two other articles
were removed because they were duplicates.
Thus, a total of 315 articles were screened. In the screening
process, the authors examined the title and the abstract for the
content and semantic relationship to the research question of
this paper. To avoid bias in the selection process, all identified
papers were individually reviewed by two reviewers.
First, all papers were rechecked for terms contained in
the search string. After this first check, the title and abstract
were checked using the exclusion criteria controlling the
semantic context of the respective work. These papers passing
this check were downloaded and stored for information
extraction.
After review by the authors, the reliability of the results
was 98.41% (agreement on 310 of 315 total articles).
No automation tools were used for the search.
As part of the screening process, several articles were
excluded to ensure the relevance and alignment of the
selected literature with the objectives of this study. A total
of 183 articles were deemed not relevant to this work for the
following reasons:
• Superficial or casual references to hallucinations:
Many articles mentioned hallucinations only briefly as a
general risk, theoretical problem, or casual observation
without providing empirical analysis, proposing new
detection or remediation methods, or providing insights
into hallucination-specific behavior in LLMs.
• Application without analytical depth: A considerable
proportion of the excluded papers implemented LLMs
for domain-specific applications without analyzing or
addressing hallucinatory behavior. These studies often
used hallucinations as a motivation, but did not directly
investigate the phenomenon or propose solutions.
• Lack of methodological or conceptual contribution:
In order to maintain the scientific rigor of the analysis,
publications were excluded if they offered little that was
new (e.g., repetition of known results), contained no
empirical evaluation, or showed minimal theoretical or
methodological development.
VOLUME 13, 2025 148235
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
• Focus on neighboring but different concerns: Several
papers addressed topics such as AI bias, trustworthi-
ness, explainability, user acceptance, or ethical use
of generative models. These studies were excluded
if hallucinations were not a central and analytically
investigated component of the work.
In total, 132 articles were included in the preliminary work
and used for the SLR (IEEE Xplore (n=76), ACM (n=44),
ScienceDirect (n=12)).
C. CLASSIFYING THE RESULTS WITH A DESCRIPTIVE
ORGANIZATIONAL SCHEME FOR LLM OUTPUTS
Of the 125 studies identified, 7 were categorized as literature
reviews and excluded, we implemented a structured classifi-
cation scheme to better understand the landscape of research
on LLM hallucinations. To this extent, we developed and
employed a descriptive organizational scheme designed to be
functionally precise and engineering-oriented. This approach
is essential for analyzing system-level failures and mitigation
strategies, as it provides a taxonomy that maps directly to
system architecture and observable behaviors. To this end,
we developed and employed a descriptive organizational
scheme to categorize the phenomena reported in the selected
literature.
This scheme is derived from a synthesis of current
research on AI-generated hallucinations and their mitigation
strategies [9]. It organizes system architecture components
and their outputs into two primary categories based on
architectural design: (A) the output of an unaugmented
LLM operating on its internal parametric knowledge, and
(B) the output of an LLM augmented with an external
DMS. A prevalent example of a DMS in the literature is
the RAG architecture [12]. This architectural distinction is
fundamental, as the introduction of a DMS creates a new
set of potential success and failure modes that must be
analyzed independently [25]. The scheme outcome details
seven mutually exclusive scenarios, which are presented
in Tables 2 through 8 and will be used to classify the
contributions of the reviewed papers.
1) SYSTEM CHARACTERISTICS WITHIN THE
ORGANIZATIONAL SCHEME
To understand the seven scenarios of output generation, it is
first necessary to define the fundamental characteristics of the
two primary system architectures identified in our scheme:
The unaugmented LLM and the LLM augmented with a
DMS.
Unaugmented LLM Characteristics: An unaugmented
LLM generates output based solely on its internal, parametric
knowledge learned during training. This architecture is
defined by the following characteristics [9], [12]:
• Reliance on parametric memory: The model’s
responses are generated directly from the knowledge
stored in the parameters of the neural network. The
behavior of the model is a direct function of the
training goal, which is usually the prediction of the next
token [12], [26].
• Static knowledge base: The internal knowledge of the
model is fixed after the training phase, which leads
to a ‘‘knowledge cutoff’’. It cannot access real-time
information or easily update its knowledge base without
being retrained or fine-tuned [9].
• Vulnerability to artifacts in the training data: The
model is prone to generating hallucinations due to
noise or bias in the rich training data. The default
training target does not explicitly reward fidelity to a
particular source document, which can lead to unfaithful
or nonsensical results [9], [11], [26].
• Single-stage generation: The output is usually generated
in a single, continuous forward pass, which makes
the process fast but relatively opaque, as there are no
intermediate components that need to be verified for
their origin [12].
Augmented LLM Characteristics: An augmented system
integrates the core LLM with one or more external compo-
nents designed to detect and mitigate hallucinations. This
architecture has a different set of operational characteristics:
• Multi-Stage Processing: These systems operate in
multiple stages, such as retrieval-then-generation (as in
RAG) or generation-then-correction. This architecture
allows for explicit points of intervention and verifica-
tion [9], [12].
• Access to Non-Parametric Memory: Many DMS
architectures, particularly RAG, augment the LLM by
providing it with access to an external, non-parametric
knowledge base (e.g., a vector index of Wikipedia). This
allows the model to ground its output in up-to-date,
verifiable information [12].
• Explicit Verification Mechanisms: The DMS com-
ponent is designed to perform fact-checking or
consistency-checking. This can involve comparing
generated text against retrieved documents using
Natural Language Inference (NLI) models or Question
Answering (QA) systems [9].
• Improved Interpretability and Controllability: By
separating knowledge retrieval from generation, these
systems can provide better interpretability. The retrieved
documents serve as proof of origin for the generated
output and explain why the model generated a particular
response [12]
• New Potential Failure Points: The added complexity
introduces new vulnerabilities. The system can fail if
the DMS retrieves irrelevant or flawed information,
or if the detection logic itself is unable to identify
subtle inconsistencies, a significant focus of current
research [9], [23].
2) SCENARIOS FOR UNAUGMENTED LLM OUTPUTS
The first two scenarios are the outcome of the baseline
behavior of a standalone LLM, as conceptually illustrated in
148236 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
FIGURE 3. Conceptual illustration of an unaugmented LLM generation
architecture.
Figure 3. Its outputs are a direct function of its training and
internal knowledge [26].
• Scenario 1: Unassisted Correct Generation represents
an ideal case where the model’s probabilistic generation
aligns with factual reality.
• Scenario 2: Unchecked Hallucination is the canonical
failure mode that has motivated a vast body of research,
an issue that has been thoroughly documented in recent
surveys [9]. Hallucinations are formally defined as
generated content that is nonsensical or unfaithful to
the provided source content [9], [10]. They are known
to arise from several sources: flaws in the training
data can cause the model to learn and store factually
incorrect information or ‘‘wrong correlations’’ within its
parameters; they can stem from parametric knowledge
gaps where the model’s knowledge is outdated or
incomplete (i.e., knowledge cutoffs); or they can be
the result of flawed decoding strategies during the
generation process itself [9]. This scenario represents
the raw, unmitigated output of a system lacking any
corrective mechanisms.
FIGURE 4. Conceptual illustration of a DMS augmented LLM architecture.
• The remaining scenarios delineate distinct and critical
failure modes within these more complex augmented
systems. Scenario 5 (Unrecognized Hallucination)
occurs when the DMS fails to detect the initial error.
[25]. This is a significant focus in the literature, often
attributed to specific failure points within the DMS
itself, such as the retrieval of irrelevant or misleading
documents, or the failure of the ‘‘detection’’ mechanism
to identify subtle inconsistencies [9], [23], [24].
• Scenario 6 (Partially Corrected Hallucination)
describes a gradual failure where the DMS addresses
some, but not all, inaccuracies in a complex output. This
can result from the failure of the LLM to fully synthesize
multiple pieces of retrieved context or from incomplete
evidence retrieval by the DMS [9].
• Finally, Scenario 7 (Corrupted by Correction) cap-
tures a paradoxical and highly problematic failure in
which the DMS itself introduces an error into an
otherwise correct LLM output. This can occur if the
retrieved knowledge is flawed or if the DMS’s corrective
logic is too aggressive, corrupting a nuanced and
accurate initial response [9]. This highlights that the
mitigation layer itself is a potential source of error that
requires independent validation.
3) SCENARIOS FOR DMS AUGMENTED LLM OUTPUTS
The other five scenarios encompass systems where a second
component is integrated to verify and/or correct the LLM’s
output, as conceptually illustrated in Figure 4. This introduces
a multi-stage process where new outcomes are possible.
• Scenario 3 (Verified Correct Generation) and Sce-
nario 4 (Successful Hallucination Correction) repre-
sent the intended successes of an augmented architec-
ture. In Scenario 4, DMS effectively converts a system
failure into a success, demonstrating the primary value
proposition of architectures like RAG, which aim to
ground LLM outputs in verifiable, retrieved evidence,
thus reducing hallucination [12].
D. ANALYSIS AND CATEGORIZATION FRAMEWORK
To systematically evaluate and categorize existing research
on hallucinations in LLMs, this study introduces an analytical
framework based on system architecture. The objective is
VOLUME 13, 2025 148237
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
to assess each research paper by examining how effectively
its proposed architecture integrates the LLM with the
DMS. Additionally, this framework explicitly considers the
seven scenarios, enhancing theoretical rigor and practical
applicability.
This framework is structured along two key dimensions:
System Architecture Categories: This dimension catego-
rizes all identified systems along four distinct architectural
archetypes. Each category is directly linked to the identified
scenarios.
• Unaugmented LLM Generation Architectures: This
category represents systems where the LLM operates
autonomously, without any structured oversight or
integrated DMS. Consequently, LLM outputs depend
exclusively on internal parametric knowledge acquired
during training. Outputs can either be factually correct
(scenario 1) or incorrect due to inherent limitations, such
as training data biases or spontaneous hallucinations
(scenario 2).
• Post-hoc Reactive Validation Architectures: This
category encompasses systems where a standalone DMS
is coupled with the LLM but operates exclusively in a
reactive manner. The DMS is triggered only after the
LLM has completed its generation, typically based on
predefined heuristics, confidence thresholds, or anomaly
signals. As a result, this architecture allows for a range
of outcomes, including Scenario 3–7, depending on
whether and how the DMS intervenes. Corrections,
if applied, are handled entirely outside the LLM, which
remains a static, unrefined generator. Evaluation criteria
for this architecture include the precision and recall of
post-hoc detection mechanisms.
• Proactive DMS Architectures: This category is
defined by a proactive and tightly coupled interaction
between the LLM and the DMS. Unlike reactive
systems, the DMS in this architecture anticipates
potential errors and actively guides or constrains the
generation process before and during output production.
This enables outcomes aligned with Scenarios 3–
7, as each output undergoes structured verification
or correction in real time. Scenarios 1 and 2 are
inherently excluded, since no output is provided
without undergoing some form of proactive validation.
Evaluation criteria for this architecture focus on the
effectiveness of preemptive error detection mechanisms
and the degree to which they suppress complex failure
modes, particularly Scenarios 5–7. These systems
often incorporate bidirectional feedback loops, enabling
dynamic refinement of the LLM’s generation pathway
based on real-time assessments of factual accuracy or
coherence.
• Integrated LLM-DMS Architectures: This category
is characterized by architectures that move beyond the
conventional separation of LLM and DMS components.
Instead, the system fully integrates both components
into a unified, dynamic, and proactive framework to
minimize hallucinations. The DMS not only anticipates
errors before but also operates through internal feedback
loops, meta-reasoning, and self-validation.
2. Domain: Classification of research based on the area
in which their approach is applied. This helps identify areas
that are most advanced in the mitigation of hallucinations
and serves as a basis for subsequent analysis. The following
domains have been identified:
• Software Engineering & Code Generation: Involves
models generating or assisting with source code, high-
lighting precision and accuracy in algorithmic logic and
syntax.
• Question Answering & Information Retrieval: Cov-
ers methods that retrieve and generate accurate and
relevant information in response to user queries, includ-
ing RAG pipelines and grounding mechanisms to reduce
factual errors.
• Ethical & Trustworthy AI: Focuses on developing
systems that behave ethically, transparently, and reliably,
addressing hallucinations through mechanisms like
toxicity detection, fact-checking, or content verification
to ensure safe deployment.
• Multimodal & Vision-Language Models: Covers
models that integrate visual, textual, or multimodal
data (e.g., image captioning, spatial reasoning, video-
language alignment), where hallucinations can arise
from cross-modal misalignment or grounding failures.
• Domain-Specific QA (e.g., Manufacturing, Legal
AI): Targets highly specialized domains (e.g., manu-
facturing, legal, financial, or risk assessment), where
factual precision is critical and hallucinations must be
addressed using domain-specific knowledge sources and
constraints.
• Biomedical & Scientific Applications: Covers sensi-
tive, high-stakes fields (e.g., medicine, biology, scien-
tific research) where factual errors can have critical con-
sequences, and hallucination mitigation often involves
rigorous domain-specific verification and citation strate-
gies.
• User-Centric AI & Explainability: Focuses on enhanc-
ing user trust through explainable, transparent outputs,
often using UI features, feedback loops, or user-guided
validation mechanisms to help users detect or correct
hallucinations.
• Autonomous Systems & Robotics: Refers to real-time
decision-making scenarios in robots and autonomous
platforms, demanding consistently reliable outputs.
• Educational AI & Learning Systems: Involves per-
sonalized learning and knowledge assessment systems,
emphasizing accuracy to ensure effective educational
outcomes.
IV. RESULTS
The SLR initially identified 132 publications related to
hallucinations in LLMs. Out of these, 125 were included
in the final analysis, while the remaining 7 were classified
148238 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
TABLE 2. Unassisted correct generation. An LLM outputs a factually correct response using only its internal knowledge, without external verification.
TABLE 3. Unchecked hallucination. A LLM produces a confident but incorrect output, with no mechanism to detect or prevent the error.
TABLE 4. Verified CORRECT GEneration. The DMS confirms the LLM’s correct output, ensuring reliability and trust.
TABLE 5. Successful hallucination correction. The DMS identifies and corrects a faulty LLM output (e.g., via retrieval), highlighting the strength of
augmented systems.
TABLE 6. Unrecognized hallucination. The LLM generates an incorrect output, and the DMS fails to intervene, creating a false sense of reliability.
TABLE 7. Partially corrected hallucination. The DMS partially corrects a flawed output, yielding a seemingly reliable response that still harbors factual
errors, posing a deceptive risk.
VOLUME 13, 2025 148239
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
TABLE 8. Corrupted by correction. The LLM generates a factually correct output, but the DMS unintentionally introduces an error, corrupting the original
response.
B. DOMAIN DISTRIBUTION OF RESEARCH
Figure 6 visualizes the distribution of the research papers
across the nine identified domains. The most prevalent
domain is Question Answering & Information Retrieval, with
32 publications, followed by Multimodal & Vision-Language
Models (28), and Domain-Specific QA (14). Conversely,
domains such as Autonomous Systems & Robotics (6)
and Educational AI & Learning Systems (5) are currently
underrepresented in the literature.
This distribution indicates that hallucinations in LLMs are
predominantly studied in high-impact application areas such
as QA and vision-language tasks.
C. DISTRIBUTION BY ARCHITECTURE CATEGORY
The papers were also analyzed according to their alignment
with the four architecture categories. Figure 7 shows the
distribution:
• Unaugmented LLM Generation Architectures: 5 papers
(4.0%)
• Post-hoc Reactive Validation Architectures: 55 papers
(44.0%)
FIGURE 5. Yearly distribution of publications included in the analysis
(blue), with excluded literature reviews shown for reference (orange).
Only papers published between 2015 and January 17, 2025, were included.
• Proactive DMS Architectures: 54 papers (43.2%)
• Integrated LLM-DMS Architectures: 11 papers (8.8%)
The majority of papers fall under the Post-hoc Reactive
as literature review papers. These reviews, although relevant
Validation Architectures and Proactive DMS Architectures
for contextual understanding, were excluded from further
categories, indicating that while many researchers explicitly
categorization and analysis to maintain focus on empirical
incorporate reasoning-based frameworks, few push the
and technical contributions.
boundary towards advanced, fully integrated systems.
The bar chart in Figure 5 presents the distribution of
the included studies by year. A clear increase in research
interest is observable in 2024, which accounts for the vast
majority of contributions (92 out of 125), while 2023 and
2025 contributed 22 and 11 studies, respectively. The sharp
rise in 2024 aligns with the growing adoption and scrutiny
of generative LLMs. The literature reviews were distributed
across 2024 (4) and 2025 (3), with none in 2023.
A. CATEGORIZATION BY DOMAIN AND SYSTEM
ARCHITECTURE
The classification framework used for this review includes
two main axes:
The classification matrix in Table 9 displays the inter-
section of these two dimensions across all 125 included
papers. This matrix serves as the foundation for the following
domain- and system architecture.
D. RESULTS BY SYSTEM ARCHITECTURE CATEGORY
1) HIGHLIGHTING UNAUGMENTED LLM GENERATION
ARCHITECTURES
This first category, Unaugmented LLM Generation Architec-
tures, is foundational to hallucination research as it includes
systems where the LLM operates autonomously, relying
exclusively on its internal parametric knowledge without
any structured oversight or integrated DMS, as illustrated
in Figure 8. Research in this area is critical for quantifying
the baseline performance and inherent limitations of LLMs,
thereby establishing the scope of the hallucination problem
that more advanced architectures seek to solve.
The studies in this category analyze the raw, uncorrected
outputs of LLMs to understand their factual accuracy and
biases. For instance, Zuccon et al. [34] directly evaluated
ChatGPT’s ability to attribute its answers to credible sources
148240 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
TABLE 9. Overview of the publications in the literature base. The table groups papers by their domain of application and categorizes them according to
the architecture exhibited in their hallucination detection and mitigation strategy. In the Domain column (1st X-axis), research is classified into categories
such as Software Engineering & Code Generation, Question Answering & Information Retrieval, Multimodal Models, and others. In the architecture
category columns (2nd - 5th X-axis), each publication is placed into one of the four previously established categories: Unaugmented LLM Generation
Architectures, Post-hoc Reactive Validation Architectures, Proactive DMS Architectures, and Integrated LLM-DMS Architectures. Each cell lists the authors
whose work falls under the respective domain and architecture classification.
and found that the model frequently hallucinates references;
a staggering 86% of the citations provided by the model did
not actually exist. Similarly, Moayeri et al. [33] developed
WorldBench to probe the factual recall of 20 state-of-the-
art LLMs, discovering significant geographic disparities in
their performance. Their analysis revealed that error rates for
factual questions about countries in Sub-Saharan Africa were
1.5 times higher than for those in North America, highlighting
systemic biases in the models’ internal knowledge bases.
Other works in this category extend the analysis to different
facets of model behavior and modality. Kraft et al. [65] inves-
tigated whether internal model enhancements, specifically
knowledge enhancement, can inherently reduce bias. Their
empirical evaluation demonstrated that such augmentations
do not necessarily prevent stereotypical associations, suggest-
ing that, without external corrective systems, even enhanced
models can perpetuate epistemic injustice. In the multimodal
domain, Jiang et al. [77] analyzed the representation space
of an unaugmented LVLM, identifying a significant gap
between textual and visual representations as a key factor
VOLUME 13, 2025 148241
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
FIGURE 8. Conceptual illustration of an unaugmented LLM generation
architecture.
responses. Their work reveals that mitigating hallucinations
in more constrained settings (e.g., yes/no questions) does
not correlate with a reduction in these more spontaneous,
unprompted fabrications, underscoring the need for bench-
marks that evaluate LLMs in their unconstrained, generative
state.
FIGURE 6. Distribution of the research papers across the nine identified
domains.
Each of these studies examines models in the absence
of external grounding or structured correction mechanisms,
making them prototypical examples of Unaugmented LLM
Generation Architectures as defined by this framework.
FIGURE 7. Distribution of the research papers across the four system
architecture categories.
contributing to hallucinations. This foundational analysis
motivated their subsequent development of a new training
methodology.
Finally, Kaul et al. [78] introduced the THRONE bench-
mark to assess hallucinations in open-ended, free-form
2) HIGHLIGHTING POST-HOC REACTIVE VALIDATION
ARCHITECTURES
Post-hoc reactive validation architectures address LLM hal-
lucinations by employing standalone detection and mitigation
systems that operate after the LLM has generated its output,
as illustrated in Figure 9. These systems typically react to
outputs based on confidence thresholds [130], heuristics [35],
[36], [89], or anomaly signals [45], [67], with corrections
occurring externally without modifying the core LLM.
Evaluation often emphasizes the precision and recall of
detection.
A common strategy across various domains involves RAG
to ground LLM outputs in external knowledge. In Question
Answering & Information Retrieval, RelD is proposed as
a robust hallucination detection discriminator trained on
the RelQA dataset, leveraging ELECTRA as a backbone
to align with human ratings [35]. FaCTQA detects and
localizes factual errors in summaries by generating QA pairs
from keywords and semantically matching answers from the
summary and source text. It outperforms baselines, achieving
90.3% accuracy and 0.8701 F1 Score on XSum [36].
Similarly, in Domain-Specific QA, Complaint-LLM uses a
knowledge-graph-driven approach with a KGCN classifier to
enhance factuality in complaint processing, achieving peak
accuracy/F1 scores of 0.85/0.86 [105]. ReRag optimizes
vector database hyperparameters in a closed feedback loop,
improving average context similarity scores from 0.494
148242 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
(before optimization) to 0.616 (after optimization) [40].
For multimodal models, RAG combined with a Visual
Hallucination Corrector in MM-PEAR-CoT enhances Visual
Spatial Description tasks, improving SPICE scores to 86.61
[86]. NewsGPT, a 70B-parameter LLM, integrates a RAG
module for hallucination control in news generation, achiev-
ing top performance in sentence-level discrimination (52.4%
accuracy) [110].
Another approach focuses on leveraging internal model
artifacts or self-assessment mechanisms for detection.
In Software Engineering & Code Generation, empirical
evaluation shows ChatGPT’s self-verification is unreliable,
with direct question prompts yielding high precision but very
low recall (e.g., 0.02 for code generation). Guiding question
prompts improve recall (+69% in vulnerability detection)
but increase false positives. Self-contradictory hallucinations
are frequently observed, especially in complex languages
like Scala [27]. For general LLMs, uncertainty estimation
techniques are explored; sample-based methods like Sample
VRO show strong negative correlations (up to−0.739)
with erroneous outputs, aiding in error detection without
internal model access [69]. Snyder et al. [45] demonstrate
that early detection of factual hallucinations is possible by
analyzing internal model artifacts like self-attention and FC
activations, achieving AUROC > 0.80. Huo et al. [39] explore
LLM self-detection of hallucinations in open-domain QA
using retrieval-based techniques, finding that post-generation
evidence retrieval can identify and correct hallucinations with
moderate success.
Furthermore, several frameworks integrate rule-based or
logic-based validation. Drowzee, a logic-programming-aided
metamorphic testing framework, detects fact-conflicting
hallucinations by deriving new facts through logic rules
and comparing LLM output graphs to ground-truth graphs.
Hallucination rates ranged from 24.7% (GPT-4) to 59.8%
(Mistral-7B) [71]. For embodied agents in Autonomous
Systems & Robotics, a two-stage alignment framework
mitigates action hallucinations by combining parameter-
efficient fine-tuning with RAG, achieving 100% Language
Compliance [142]. CoT-TL translates natural language
planning instructions into Linear Temporal Logic (LTL)
formulas, integrating a model checker to reject invalid
LTL formulas and self-consistency for robustness, achieving
up to 91.69% accuracy [143]. In Biomedical & Scientific
Applications, IC-BERT, a BERT-based instruction classifier,
pre-screens user input to reduce hallucinations in TCM
queries, achieving 99.95% accuracy [119]. User-centric AI
initiatives like HILL aim to help users detect hallucinations
through interface features like confidence scores and source
links, showing increased hallucination awareness (+83.4%)
in user studies [130].
In each case, the LLM remains a ‘‘static, unrefined
generator’’; validation is delegated to a standalone module
that fires reactively and either blocks unsafe content or
flags and corrects it after the fact. The lack of in-flight
reasoning correction, coupled with an exclusive focus on
FIGURE 9. Conceptual illustration of a Post-hoc reactive validation
architecture.
precision/recall of these downstream detectors, is precisely
why this work belongs in the subsection ‘‘Post-hoc Reactive
Validation Architectures’’.
3) HIGHLIGHTING PROACTIVE DMS ARCHITECTURES
Research in this category shifts from post-hoc validation to
proactive intervention, where the DMS is tightly coupled
with the LLM to guide and constrain the generation
process in real time. These architectures are defined by
their ability to anticipate and preemptively correct errors,
often through multi-stage pipelines, formal verification, and
iterative feedback loops.
A dominant approach involves enhancing RAG with
verification stages. In Question Answering, frameworks
decompose complex queries into a series of sub-questions,
where each reasoning step must be grounded in retrieved
evidence before the final answer is synthesized [53], [56].
Some systems embed fact-checking directly into the retriever
pipeline, pruning irrelevant or non-factual documents before
they can influence the generator [48], or use real-time
analysis of token entropy to trigger retrieval only at points
of high uncertainty, improving efficiency [57]. This principle
extends to domain-specific applications, where structured
data from financial reports or manufacturing glossaries is
converted into hierarchical chunks, and queries are enhanced
with verified terminologies before LLM inference [111],
[112], [114].
In high-stakes domains like Software Engineering & Code
Generation, proactive systems integrate formal methods to
enforce correctness. For instance, Refine4LLM uses a theo-
rem prover to formally verify each refinement step suggested
VOLUME 13, 2025 148243
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
by the LLM, creating a robust generation-validation loop
that ensures code is provably correct [29]. Similarly,
other frameworks employ multi-phase filtration pipelines
that guarantee generated unit tests successfully build, pass
reliably, and increase code coverage. At Meta, TestGen-
LLM achieved a 75% build success rate, 57% pass rate,
and 25% of the test cases increased coverage, significantly
reducing hallucinations through formal guarantees [28].
This verification-in-the-loop pattern is also applied in
Autonomous Systems, where formal verifiers check the
LLM’s plans and feed counterexamples back into the prompt
for iterative correction. Preliminary experiments show that
this approach can successfully eliminate identified halluci-
nations through iterative refinement using counterexamples,
thereby converging toward verified outputs [145].
Multimodal systems employ novel decoding-time strate-
gies that actively reshape the generation path. Techniques like
DOPRA and OPERA monitor internal self-attention maps
during generation, penalizing or rolling back the process
when patterns indicative of ‘‘over-trust’’ in summary tokens
emerge, thereby preventing the model from deviating from
visual facts [95], [100]. Other frameworks use fine-grained
correctional human feedback on specific hallucinated seg-
ments to retrain the model with a loss function that heavily
prioritizes corrected regions. For instance, RLHF-V reduced
the object hallucination rate by 34.8%, surpassing baselines
trained with 7×more data [76]. Despite these advances,
failure modes persist; for example, formal verifiers can time
out on complex specifications (Scenario 5), and RAG systems
can be corrupted by noisy or outdated retrieved data (Scenario
7), highlighting that even proactive systems require careful
validation of the DMS itself [29], [116].
Ultimately, these architectures signify a fundamental shift,
treating the LLM not as a static generator to be checked but as
an active component in a dynamic reasoning and verification
loop. The emphasis on in-flight correction and preemptive
guidance is precisely what defines these systems as proactive.
FIGURE 10. Conceptual illustration of proactive DMS architectures.
4) HIGHLIGHTING INTEGRATED LLM-DMS ARCHITECTURES
Integrated LLM-DMS architectures, as illustrated in
Figure 11, mark the current frontier in hallucination
mitigation by dissolving the traditional boundary between
generator and detector. In these systems, the DMS is woven
directly into the reasoning loop of the language model and
uses internal feedback, meta-reasoning, and self-validation
to anticipate, intercept, and correct errors before they occur.
Early examples include the retrieval-optimization-validation
chain of Shi et al. [62], in which a self-critique module
continuously refines the retriever’s reasoning chain, and
FS-C of Sun et al. [63], which enriches this loop with
multi-perspective filtering to block faulty evidence in flight.
The CORE framework by Wang et al. [64] adds causal
graphs and Pearl-style interventions so that scoring and
interventional regeneration jointly increase factual reliability.
Experiments showed that CORE improved accuracy on
HotpotQA by 8.4% and reduced hallucination by 11.3% on
ScienceQA compared to state-of-the-art baselines.
Other work focuses on linking intermediate inferences
to explicit memory structures. The FACT framework by
Gao et al. [102] couples rationale generation with post-hoc
faithfulness checks and feeds validated rationales back into
inference, while the GROUNDHOG system by Zhang et al.
[104] ties language reasoning to symbolic visual mem-
ory so that a holistic reasoner can critique and revise
multimodal outputs. GROUNDHOG achieves 22.6% lower
object hallucination rates on Grounded Image Captioning
(GIC) tasks and outperforms prior state-of-the-art methods
without any task-specific fine-tuning. Complementing these
architectures, the LVLM-eHub benchmark by Xu et al. [103]
provides systematic, human-annotated probes for object
hallucinations, supporting the co-development of integrated
vision–language models.
Domain-specific integrations bring additional guarantees.
In combustion science, Sharma and Raman [127] layer
reference verification, physics-guided checks, and structured
fallback generation onto a RAG backbone, maintaining
scientific fidelity across literature review and simulation
planning. MedLaSA [128] applies causal tracing to edit
only the neurons responsible for a faulty fact, improving
medical accuracy without disrupting unrelated knowledge.
Sun et al. [139] their adaptive in-context learning paradigm
and Zhang et al. [140] their empathetic conversational rec-
ommender both add auto-feedback and sentiment-aware
reweighting to keep personalised outputs consistent with user
intent and emotion, whereas Shi et al. [151] their multi-agent
collaborative filtering framework distributes validation across
interacting agents, reducing residual errors in the most
complex Scenario 5-7 settings. Specifically, the MCF method
improved average accuracy across ten datasets by 9.1% over
148244 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
FIGURE 11. Conceptual illustration of an integrated LLM-DMS
architecture.
baseline and achieved 75% valid reasoning trajectories in
symbolic tasks.
The defining feature of these studies is a closed, antic-
ipatory control loop. The DMS interrogates intermediate
reasoning, intervenes before an error is finalized, and
self-improves when its own judgments prove unreliable.
By unifying generation, detection, and correction in a single
dynamic pipeline, integrated LLM-DMS architectures set a
new benchmark for long-term factual robustness.
V. DISCUSSION
A. KEY CHALLENGES AND OPEN ISSUES
One notable finding of our study is the uneven distribution
of domains in hallucination research. Two areas - Question
Answering & Information Retrieval and Multimodal &
Vision-Language Models - account for 48% of all contri-
butions, while the remaining seven areas account for only
52%, indicating a strong bias. This is particularly noticeable
in areas such as Software Engineering & Code Generation
(only 6 contributions, 4.8%), which is surprising given the
widespread promotion of LLMs as a tool for programming
tasks. Also underrepresented are the areas of Educational
AI & Learning Systems, as well as Autonomous Systems
& Robotics, which are only researched to a limited extent
despite their growing social importance.
The discrepancy is even clearer when looking at the
architectural sophistication of the systems. No contribution
from the area of Software Engineering & Code Generation
was rated above the Proactive DMS Architectures category,
whereas in the areas of Question Answering & Information
Retrieval [62], [64], [139] as well as in the areas of
Multimodal & Visual-Language Models [102], [103], [104],
three contributions each exceeded this threshold and thus
demonstrated more advanced and integrated approaches to
detecting and mitigating hallucinations.
In high-stakes domains like Biomedical & Scientific
Applications, the need for hallucination control is critical.
Encouragingly, research by Sharma and Raman [127] and
Xu et al. [128] presents techniques that surpass traditional
reactive correction models, indicating the development of
robust frameworks that significantly reduce hallucinations,
even extending to fully Integrated LLM-DMS Architectures.
On the other hand, domains such as Ethical & Trustwor-
thy AI, while not underrepresented in terms of quantity,
show no contributions rated at the Integrated LLM-DMS
Architectures category. Although two notable efforts [75],
[76] achieved Proactive DMS Architectures, most papers
remain at a Post-hoc Reactive Validation Architectures level,
indicating an opportunity for improvement in implementing
more advanced and integrated mitigation strategies.
Other areas, such as Domain-Specific QA, show a high
level of architectural sophistication in many papers (eight
papers rated Proactive DMS Architectures), but still lack Inte-
grated LLM-DMS Architectures, revealing a ripe opportunity
for further advances. Similarly, Autonomous Systems &
Robotics are not only underrepresented, but also do not show
a mature application of advanced architectural frameworks,
highlighting the need for targeted future research.
B. THEORETICAL IMPLICATIONS
The analysis of techniques related to hallucinations based on
system architecture has shown that many current approaches
remain heuristic or dominated by single-stage processing,
particularly in less-researched areas. The utility of the
framework lies in highlighting which domains have evolved
towards deliberative, evaluative, and multi-phase processes
(Proactive DMS Architectures and Integrated LLM-DMS
Architectures) and which continue to rely on reactive or
single-layer solutions.
The fact that only 11 papers (8.8%) achieved Integrated
LLM-DMS Architectures reflects a greater challenge in
aligning machine reasoning with highly integrated and
sophisticated architectural designs. It raises the question of
whether more deeply integrated architectural modeling is
needed to advance hallucination mitigation beyond superfi-
cial corrections.
Moreover, success in high-stakes domains such as
Biomedicine & Scientific Applications suggests that risk
awareness can drive architectural sophistication in sys-
tem development. This insight can help theorists develop
context-sensitive architectural extensions that are better
suited for domain-specific challenges.
C. PRACTICAL IMPLICATIONS
From a practical perspective, this review provides a diagnos-
tic overview of where hallucination research is concentrated,
and where it is not. Practitioners working in underrepresented
domains such as Software Engineering & Code Generation,
Educational AI & Learning Systems, and Autonomous Sys-
tems & Robotics should be aware of the methodological gaps
and consider adapting techniques from more architecturally
VOLUME 13, 2025 148245
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
sophisticated domains. The lack of Integrated LLM-DMS
Architectures in these areas suggests that current systems
may not be adequately equipped to detect and/or prevent
hallucinations in high-impact scenarios.
In contrast, the architectural sophistication seen in
Biomedical & Scientific Applications domain shows that
multi-phase, knowledge-grounded, and uncertainty-aware
architectures are both feasible and beneficial [127], [128].
Developers and product teams may draw from these domains
to implement cross-domain mitigation techniques or to
design hybrid models that combine strengths across research
areas.
Finally, the taxonomy and mapping presented in this
paper serve as a reference point for both researchers and
engineers. It provides a foundation for benchmarking future
techniques, understanding their architectural characteristics
and integration levels, and identifying transferable practices.
D. CROSS-DOMAIN TRANSFERABILITY OF INTEGRATED
LLM-DMS ARCHITECTURES
If we assume that Integrated LLM-DMS Architectures repre-
sent a real advance in mitigating hallucinations by enabling
systems to anticipate, recognize, and correct errors through
reasoning and meta-reasoning mechanisms, then a crucial
next question arises: Can these advanced architectural
approaches be meaningfully adapted to domains beyond
those for which they were originally developed?
Empirical support for such cross-domain portability is
growing. For instance, in the ontology-engineering study
by Doumanas et al. [14], a four-level human-to-LLM
collaboration scheme was applied to both medical and
search-and-rescue ontologies, demonstrating that staged,
human-in-the-loop designs remain effective when ported
across knowledge domains. Further evidence comes from
research on high-fidelity training datasets. Dziri et al. [152]
created FAITHDIAL, a benchmark of hallucination-free
dialogues, and found that models trained on this curated data
showed remarkable generalization. In a zero-shot transfer
setting, their model significantly reduced hallucinations
when tested on other knowledge-grounded dialogue datasets
like TopicalChat and CMU-DoG, proving more faithful
than models trained specifically on the in-domain data.
Similarly, the successful application of general-purpose
architectures to highly specialized fields highlights this
transferability. Lee et al. [153] demonstrated that a RAG
model, a common strategy for mitigating hallucination,
outperformed general-purpose models like GPT-4 in the
specialized domain of construction safety management. Their
work underscores how architectural solutions designed for
broad knowledge-intensive tasks can be effectively adapted
to improve factual grounding and reliability in niche, high-
stakes domains.
While approaches categorized as Integrated LLM-DMS
Architectures show promising results in domains such as
Biomedicine, Multimodal AI, and Educational Systems,
their architectures often remain tightly coupled to their
domain-specific assumptions. This raises the question of
whether and how they can be transferred to underrepresented
yet practically important areas such as Software Engineering
& Code Generation, Ethical & Trustworthy AI, Domain-
Specific Questioning & Answering, and Autonomous Sys-
tems & Robotics and what adaptations, limitations, or risks
might be involved in doing so.
A comparative review of the Integrated LLM-DMS
Architecture contributions reveals several architectural char-
acteristics that are not only advanced but also potentially
transferable across domains. Notable examples include
multi-agent validation mechanisms [151], causal knowl-
edge pipelines [64], adapter-based knowledge editing with
minimal collateral effect [128], and retrieval-augmented
generation with embedded citation verification [127]. These
strategies often exhibit modularity, feedback loops, or self-
correction mechanisms that can be adapted with reasonable
effort.
Other approaches, such as rationale distillation through
symbolic execution and natural language transforma-
tion [102], sentence-weighted reasoning validation [63],
or grounding via segmentation and pointer tokens [104],
also reveal underlying principles that could be tailored
for different domain challenges. Although frameworks like
LVLM-eHub [103] serve primarily as benchmarks rather than
intervention mechanisms and are therefore excluded from
transferability, the vast majority of Integrated LLM-DMS
Architecture contributions indicate opportunities for reuse in
different system environments.
By analyzing the papers classified as Integrated LLM-DMS
Architectures in conjunction with the defined application
domains, it was possible to identify components and
architectural principles within these approaches that hold
strong potential for cross-domain transfer. Specifically, the
analysis revealed that several of these methods exhibit
modular reasoning units, meta-validation layers, or causal
grounding strategies that could be meaningfully adapted to
domains that are currently underrepresented in the literature.
The following discussion explores this transferability along
four underrepresented domains: Software Engineering &
Code Generation, Ethical & Trustworthy AI, Domain-
Specific Question Answering, and Autonomous Systems &
Robotics, highlighting the rationale for adaptation, possible
benefits, and domain-specific constraints.
1) SOFTWARE ENGINEERING & CODE GENERATION
Several Integrated LLM-DMS Architecture approaches show
strong promise for transfer into the Software Engineering
domain. For instance, the Fact framework [102], which
generates verifiable and concise rationales via symbolic
execution and natural language transformation, could sig-
nificantly benefit automated code explanation, debugging,
and unit test generation. The symbolic-to-natural logic
transformation pipeline aligns with the needs of explainable
and deterministic software behavior modeling.
148246 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
Similarly, Sharma and Raman [127] citation-verified RAG
pipeline, initially designed for combustion science, could
be adapted to structure large-scale code repositories and
deliver grounded, citation-supported answers in software
maintenance or DevOps scenarios.
2) ETHICAL & TRUSTWORTHY AI
Several frameworks, including Multi-agent Collaborative
Filtering (MCF) [151] and CORE [64], offer mechanisms that
could improve trust, explainability, and ethical compliance
in generative systems. For example, MCF’s multi-agent ver-
ification and AMR-based alignment may help in validating
ethical consistency across diverse outputs. Likewise, CORE’s
causal reasoning pipeline could be employed to trace and
justify ethically sensitive decisions.
MedLaSA [128], while tailored to biomedical knowledge,
also provides fine-grained editing and control with minimal
collateral effects, features vital for updating knowledge bases
in high-stakes ethical scenarios such as regulatory changes or
bias mitigation.
3) DOMAIN-SPECIFIC QA
The causal pipeline in CORE [64], with its use of belief
scoring and do-calculus, is directly relevant to structured QA
tasks in domains such as law, manufacturing, or finance.
Likewise, MedLaSA’s [128] scalable adapter injection could
allow fine-grained factual corrections in specialized corpora,
such as technical standards or compliance regulations.
Furthermore, the Groundhog framework [104], while orig-
inally multimodal, offers a transparent grounding mechanism
that could support visual QA tasks in technical drawings or
diagnostic imaging.
4) AUTONOMOUS SYSTEMS & ROBOTICS
Transferability into autonomous systems requires high rea-
soning robustness, low latency, and explainable decision-
making. In this context, GROUNDHOG’s segmentation-
based grounding and multi-modal validation [104] could
help mitigate hallucinations in vision-based robotic control.
Similarly, MCF’s agent-based architecture [151] aligns well
with ensemble-based decision-making in robotics.
Additionally, Fact’s rationale validation mechanism [102]
may aid in producing human-interpretable justifications for
autonomous decisions, a key requirement for human-in-the-
loop systems.
While these domain-specific transfer scenarios show
promising ways to reuse the techniques of Integrated
LLM-DMS Architectures, they are not without limitations.
Cross-domain adaptation inherently involves challenges that
must be addressed to ensure effectiveness and reliability.
In resource-poor domains, the lack of adequate ground
truth [64] could compromise the effectiveness of validation
mechanisms. Uncalibrated adaptation could also lead to
semantic drift, where correction mechanisms misrecognize
or distort domain-specific reasoning patterns [63], [102].
Furthermore, the increasing modularity and abstraction of
Integrated LLM-DMS Architectures might benefit gener-
alization but reduce transparency, leading to trust and
explanation problems in high-stakes domains [102], [104].
Nevertheless, the reviewed approaches collectively repre-
sent a rich design space for hallucination mitigation, many of
which hold potential beyond their initial context.
By fostering architecture transfer from mature to under-
represented domains and modular adoption of Integrated
LLM-DMS Architecture components, hallucination research
can be extended to broader application domains, ultimately
helping to close the gap between system reliability and real-
world deployment.
E. BRIDGING LLM ARCHITECTURES WITH DUAL PROCESS
THEORY
Beyond the architectural classification and transferability
insights discussed in this review, future research could benefit
from a deeper exploration of the cognitive underpinnings of
hallucinations in LLMs. In particular, the application of Dual
Process Theory (DPT) offers a promising complementary
framework for understanding how LLMs generate outputs,
especially in contexts prone to factual inaccuracy.
DPT, rooted in foundational work from cognitive psy-
chology and decision science [154], [155], [156], [157],
distinguishes between two qualitatively different cognitive
modes: System 1, which is fast, intuitive, and associative, and
System 2, which is slow, reflective, and rule-based. While this
review has focused primarily on system-level architectural
classifications of hallucination mitigation strategies, future
work could explore whether and how these DPT constructs
can be operationalized within the internal mechanisms of
LLMs.
Recent advances indicate that such integration is feasible.
For instance, Cheng et al. [158] introduce HaluSearch,
a framework that employs tree-search-based slow thinking
and a dynamic System 1 / System 2 switch to mitigate
hallucinations during inference. This approach maps DPT
concepts onto computational processes, replacing metaphor
with empirically testable mechanisms such as step-level,
self-evaluation and hierarchical reasoning selection. Such
implementations provide a promising direction for further
empirical investigation.
Nonetheless, existing DPT-inspired models for LLMs have
been critiqued for oversimplification or ad hoc conceptual
mapping [159]. To move beyond metaphorical alignment,
future research should aim to establish concrete empirical
mappings between neural activations or internal states in
LLMs and cognitive constructs from DPT. For example,
studies could investigate whether ‘‘slow thinking’’ modules,
such as meta-cognitive loops or deliberative planning agents,
correlate with measurable improvements in factual consis-
tency, response calibration, or self-correction behaviors.
A more granular application of DPT may also illuminate
failure modes in current systems. Whereas most current
hallucination mitigation strategies either operate globally
VOLUME 13, 2025 148247
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
(e.g., via training data quality) or post hoc (e.g., through
faithfulness scoring), a DPT-inspired approach could enable
real-time cognitive modulation of reasoning depth based
on context complexity. This opens up the possibility of
adaptive control mechanisms that switch between System
1-like generation and System 2-like reasoning based on task
demands.
Finally, future work should also investigate the boundary
conditions and limitations of DPT when applied to artificial
systems. As Evans [155] emphasizes, dual-process theories
encompass diverse and sometimes conflicting accounts,
and not all DPT attributes may sensibly map onto LLM
architectures. Moreover, as Norman et al. [157] and others
note, both intuitive and deliberative processes can contribute
to error under specific conditions. A critical research
priority will thus be to empirically identify when and how
slow-thinking mechanisms meaningfully improve reliability
without introducing unacceptable computational overhead or
unintended consequences.
In sum, by systematically investigating the relevance and
implementation of DPT-inspired constructs in LLM architec-
tures, future research can move beyond superficial analogies
and toward principled, theory-guided system design. This
may ultimately lead to the development of more robust,
explainable, and human-aligned generative models.
VI. CONCLUSION
This paper presents a SLR of hallucinations in LLMs,
which provides a domain-based classification integrated
into a descriptive organizational scheme based on system
architectures. By analyzing 125 peer-reviewed publications,
we were able to identify patterns in the way hallucination
detection and mitigation techniques are applied in nine
application domains and how they relate to different levels
of architectural sophistication and integration.
Our findings show that research efforts are concentrated
in the areas of Question Answering & Information Retrieval,
as well as Multimodal & Visual-Language Models, which
together account for almost half of all studies (48%). These
domains also exhibit higher architectural sophistication,
including several approaches that demonstrate Integrated
LLM-DMS Architectures by incorporating multiphase pro-
cessing, knowledge-based, and uncertainty-aware compo-
nents. In contrast, areas such as Software Engineering &
Code Generation, Educational AI & Learning Systems, and
Autonomous Systems & Robotics are underrepresented and
underdeveloped despite their practical relevance in terms of
hallucination resilience.
A. LIMITATIONS
While this review provides a comprehensive account of
techniques related to hallucinations, several limitations must
be acknowledged.
First, the inclusion criteria and keyword filters may have
inadvertently excluded relevant papers, particularly those that
used non-standardized terminology or focused on adjacent
tasks (e.g., summarizing or paraphrasing) without explicitly
identifying them as hallucination-related.
Second, although guided by a standardized rubric, the
assignment of system architecture categories involves a
certain degree of subjectivity, especially for borderline cases
or hybrid methods that do not clearly fit into the defined
categories. In addition, the classification of domains may
be fluid for some interdisciplinary work, which could easily
influence the number of domains.
Third, a further limitation concerns the organizational
scheme applied in this review. Specifically, we employed a
set of simplified architectural categories to classify detection
and mitigation strategies. While these abstractions provided
a structured and comprehensible lens through which to
evaluate the literature, they also simplified the rich and often
overlapping features of complex LLM-DMS interactions.
By mapping hallucination-related techniques in LLMs onto
these simplified categories, we risked compressing the
underlying algorithmic dynamics into discrete classifications
that may not fully capture their nuances. However, adopting
this simplified framework enabled a systematic comparison
across a large body of work and supported the creation of a
unified classification scheme. Future work could extend this
analysis using more granular or hybrid architectural models
that reflect new insights into machine reasoning and provide
a more in-depth explanation of how hallucinations arise and
are processed in artificial systems.
Fourth, although several promising Integrated LLM-DMS
Architecture approaches with cross-domain potential have
been identified in this review, the actual process of
transferability remains largely theoretical. We have not
experimentally tested how well these frameworks work
outside their original domains. The effectiveness, scalability,
and potential limitations of these techniques in adapting to
new domains, therefore, remain an open empirical question
that warrants future investigation.
Lastly, our analysis does not assess quantitative perfor-
mance metrics for the different papers, but focuses on
methodological features, which may limit direct comparabil-
ity in terms of effectiveness.
B. FUTURE RESEARCH
This work presents a theoretical classification of hallucina-
tion mitigation techniques using a domain-based framework
aligned with system architectures. Building on the results of
this study, several areas can be explored in future research.
One important direction is the empirical validation and
practical adaptation of Integrated LLM-DMS Architec-
ture techniques in underrepresented domains. While this
review already analyzed the transferability potential of these
advanced architectural approaches across domains such as
Software Engineering, Educational AI, and Autonomous Sys-
tems this analysis remains conceptual. What is now required
are applied studies that benchmark these architectures in
target domains, under real-world conditions. Techniques
such as multi-agent validation [151], and causal reasoning
148248 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
pipelines [150] have shown promise in a cursory analysis in
Ethical & Trustworthy AI and Domain-Specific QA domains.
Their adaptation must now be tested empirically to assess
not only effectiveness but also implementation feasibility,
performance trade-offs, and integration challenges within
new domain-specific workflows.
To support this, future work should also aim to isolate
the most impactful architectural modules through ablation
studies and module-level benchmarking. Evaluating compo-
nents such as feedback loops, meta-reasoning layers [62],
and retrieval-integrated validators [127] could guide the
design of modular toolkits for cross-domain reuse. Bridging
the gap between theory and practice will require not
only empirical validation but also the systematic devel-
opment of reusable and adaptable architectural building
blocks.
Additionally, the subjective nature of architectural sophis-
tication assessment and the variability in domain boundaries
present challenges that may be addressed in future studies
by developing more fine-grained evaluation rubrics or
automated classification tools. This could also include the
creation of a modular benchmarking framework that system-
atically compares components such as retrieval integration,
error correction, and feedback loops across tasks.
Findings from the discussion also suggest potential
in exploring hallucination behaviors in multimodal and
conversational systems, where factors like co-reference
resolution and modality fusion can introduce new types
of hallucination [93]. Finally, deeper integration of user
feedback and system alignment mechanisms, including
active learning, reinforcement learning, or interactive error
correction, may help advance the architectural sophistication
of techniques and improve the robustness of LLM-based
systems.
In summary, future research has the opportunity not
only to close domain-specific gaps but also to advance
theoretical models toward a more unified, empirical, and
context-sensitive understanding of hallucination in LLMs.
REFERENCES
[1] D. Bzdok, A. Thieme, O. Levkovskyy, P. Wren, T. Ray, and S. Reddy,
‘‘Data science opportunities of large language models for neuroscience
and biomedicine,’’ Neuron, vol. 112, no. 5, pp. 698–717, Mar. 2024.
[2] E. Hsu and K. Roberts, ‘‘Leveraging large language models for
knowledge-free weak supervision in clinical natural language process-
ing,’’ Sci. Rep., vol. 15, no. 1, Mar. 2025, Art. no. 8241.
[3] K. K. Singhal et al., ‘‘Toward expert-level medical question answering
with large language models,’’ Nature Med., vol. 31, no. 3, pp. 943–950,
2025.
[4] B. Romera-Paredes, M. Barekatain, A. Novikov, M. Balog, M. P. Kumar,
E. Dupont, F. J. R. Ruiz, J. S. Ellenberg, P. Wang, O. Fawzi, P. Kohli,
and A. Fawzi, ‘‘Mathematical discoveries from program search with large
language models,’’ Nature, vol. 625, no. 7995, pp. 468–475, Jan. 2024.
[5] A. Algherairy and M. Ahmed, ‘‘A review of dialogue systems: Current
trends and future directions,’’ Neural Comput. Appl., vol. 36, no. 12,
pp. 6325–6351, Apr. 2024.
[6] C. Zhai, S. Wibowo, and L. D. Li, ‘‘Evaluating the AI dialogue system’s
intercultural, humorous, and empathetic dimensions in English language
learning: A case study,’’ Comput. Educ., Artif. Intell., vol. 7, Dec. 2024,
Art. no. 100262.
[7] Y. Shi, P. Ren, J. Wang, B. Han, T. ValizadehAslani, F. Agbavor, Y. Zhang,
M. Hu, L. Zhao, and H. Liang, ‘‘Leveraging GPT-4 for food effect
summarization to enhance product-specific guidance development via
iterative prompting,’’ J. Biomed. Informat., vol. 148, pp. 1–9, Dec. 2023.
[8] T. Zhang, F. Ladhak, E. Durmus, P. Liang, K. McKeown, and T. B.
Hashimoto, ‘‘Benchmarking large language models for news summa-
rization,’’ Trans. Assoc. for Comput. Linguistics, vol. 12, pp. 39–57,
Jan. 2024.
[9] Z. Ji, N. Lee, R. Frieske, T. Yu, D. Su, Y. Xu, E. Ishii, Y. J. Bang,
A. Madotto, and P. Fung, ‘‘Survey of hallucination in natural language
generation,’’ ACM Comput. Surveys, vol. 55, no. 12, pp. 1–38, Dec. 2023.
[10] S. Farquhar, J. Kossen, L. Kuhn, and Y. Gal, ‘‘Detecting hallucinations
in large language models using semantic entropy,’’ Nature, vol. 630,
no. 8017, pp. 625–630, Jun. 2024.
[11] J. Maynez, S. Narayan, B. Bohnet, and R. McDonald, ‘‘On faithfulness
and factuality in abstractive summarization,’’ in Proc. 58th Annu. Meeting
Assoc. Comput. Linguistics, 2020, pp. 1906–1919.
[12] P. Lewis, E. Perez, A. Piktus, F. Petroni, V. Karpukhin, N. Goyal,
H. Küttler, M. Lewis, W.-T. Yih, T. Rocktäschel, S. Riedel, and D.
Kiela, ‘‘Retrieval-augmented generation for knowledge-intensive NLP
tasks,’’ in NIPS’22: Proc. 36th Int. Conf. Neural Inf. Process. Syst., 2020,
pp. 9459–9474.
[13] J. Lee, X. Wang, D. Schuurmans, M. Bosma, E. H., Q. V. Le,
and D. Zhou, ‘‘Chain-of-thought prompting elicits reasoning in large
language models,’’ in Proc. 36th Int. Conf. Neural Inf. Process.
Syst., 2022, pp. 24824–24837.
[14] D. Doumanas, G. Bouchouras, A. Soularidis, K. Kotis, and G. Vouros,
‘‘From human-to LLM-centered collaborative ontology engineering,’’
Appl. Ontology, vol. 19, no. 4, pp. 334–367, Nov. 2024.
[15] G. P. Reddy, Y. V. Pavan Kumar, and K. P. Prakash, ‘‘Hallucinations
in large language models (LLMs),’’ in Proc. IEEE Open Conf. Electr.,
Electron. Inf. Sci. (eStream), Apr. 2024, pp. 1–6.
[16] G. Perković, A. Drobnjak, and I. Botički, ‘‘Hallucinations in LLMs:
Understanding and addressing challenges,’’ in Proc. 47th MIPRO ICT
Electron. Conv. (MIPRO), May 2024, pp. 2084–2088.
[17] E. Lavrinovics, R. Biswas, J. Bjerva, and K. Hose, ‘‘Knowledge graphs,
large language models, and hallucinations: An NLP perspective,’’ J. Web
Semantics, vol. 85, pp. 1–7, May 2025.
[18] E. Sanu, T. K. Amudaa, P. Bhat, G. Dinesh, A. U. K. Chate, and
P. R. Kumar, ‘‘Limitations of large language models,’’ in Proc. 8th Int.
Conf. Comput. Syst. Inf. Technol. Sustain. Solutions (CSITSS), Nov. 2024,
pp. 1–6.
[19] N. Maleki, B. Padmanabhan, and K. Dutta, ‘‘AI hallucinations: A
misnomer worth clarifying,’’ in Proc. IEEE Conf. Artif. Intell. (CAI),
Jun. 2024, pp. 133–138.
[20] L. Huang, W. Yu, W. Ma, W. Zhong, Z. Feng, H. Wang, Q. Chen, W. Peng,
X. Feng, B. Qin, and T. Liu, ‘‘A survey on hallucination in large language
models: Principles, taxonomy, challenges, and open questions,’’ ACM
Trans. Inf. Syst., vol. 43, no. 2, pp. 1–54, Mar. 2025.
[21] R. Zhang, H.-W. Li, X.-Y. Qian, W.-B. Jiang, and H.-X. Chen, ‘‘On large
language models safety, security, and privacy: A survey,’’ J. Electron. Sci.
Technol., vol. 23, no. 1, Mar. 2025, Art. no. 100301.
[22] M. J. Page et al., ‘‘The PRISMA 2020 statement: An updated guideline
for reporting systematic reviews,’’ Systematic Rev., vol. 10, no. 1, pp. 1–9,
Dec. 2021.
[23] A. Wang, K. Cho, and M. Lewis, ‘‘Asking and answering questions to
evaluate the factual consistency of summaries,’’ in Proc. 58th Annu.
Meeting Assoc. Comput. Linguistics, 2020, pp. 5008–5020.
[24] P. Laban, T. Schnabel, P. N. Bennett, and M. A. Hearst, ‘‘SummaC:
Re-visiting NLI-based models for inconsistency detection in summa-
rization,’’ Trans. Assoc. for Comput. Linguistics, vol. 10, pp. 163–177,
Feb. 2022.
[25] H. Zhang, H. Deng, J. Ou, and C. Feng, ‘‘Mitigating spatial hallucination
in large language models for path planning via prompt engineering,’’ Sci.
Rep., vol. 15, no. 1, pp. 1–13, Mar. 2025.
[26] L. Ouyang, J. Wu, X. Jiang, D. Almeida, C. L. Wainwright, P. Mishkin,
C. Zhang, S. Agarwal, K. Slama, A. Ray, J. Schulman, J. Hilton, F. Kelton,
L. E. Miller, M. Simens, A. Askell, P. Welinder, P. Christiano, J. Leike,
and R. Lowe, ‘‘Training language models to follow instructions with
human feedback,’’ in Proc. 36th Int. Conf. Neural Inf. Process. Syst.,
2022, pp. 27730–27744.
VOLUME 13, 2025 148249
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
[27] X. Yu, L. Liu, X. Hu, J. W. Keung, J. Liu, and X. Xia, ‘‘Fight fire with fire:
How much can we trust ChatGPT on source code-related tasks?’’ IEEE
Trans. Softw. Eng., vol. 50, no. 12, pp. 3435–3453, Dec. 2024.
[28] N. Alshahwan, J. Chheda, A. Finogenova, B. Gokkaya, M. Harman,
I. Harper, A. Marginean, S. Sengupta, and E. Wang, ‘‘Automated unit test
improvement using large language models at meta,’’ in Companion Proc.
32nd ACM Int. Conf. Found. Softw. Eng., Jul. 2024, pp. 185–196.
[29] Y. Cai, Z. Hou, D. Sanan, X. Luan, Y. Lin, J. Sun, and J. S. Dong,
‘‘Automated program refinement: Guide and verify code large language
model with refinement calculus,’’ Proc. ACM Program. Lang., vol. 9,
no. POPL, pp. 2057–2089, Jan. 2025.
[30] H. A. Quddus, M. S. Hossain, Z. Cevahir, A. Jesser, and M. N. Amin,
‘‘Enhanced VLSI assertion generation: Conformingto high-level speci-
fications and reducing LLMHallucinations with RAG,’’ in Proc. Design
Verification Conf. Exhib. Eur., 2024, pp. 57–62.
[31] S. Samanta, O. Chatterjee, H. Gupta, P. Mohapatra, A. De Maga-
lhaes, A. Rahane, M. Palaci-Olgun, and R. Kattinakere, ‘‘Efficient
incident summarization in ITOps: Leveraging entity-based grouping,’’
in Proc. IEEE Int. Conf. Softw. Services Eng. (SSE), Jul. 2024,
pp. 97–103.
[32] M. H. Tanzil, J. Y. Khan, and G. Uddin, ‘‘ChatGPT incorrectness
detection in software reviews,’’ in Proc. IEEE/ACM 46th Int. Conf. Softw.
Eng., Apr. 2024, pp. 1–12.
[33] M. Moayeri, E. Tabassi, and S. Feizi, ‘‘WorldBench: Quantifying
geographic disparities in LLM factual recall,’’ in Proc. ACM Conf.
Fairness, Accountability, Transparency, Jun. 2024, pp. 1211–1228.
[34] G. Zuccon, B. Koopman, and R. Shaik, ‘‘ChatGPT hallucinates when
attributing answers,’’ in Proc. Annu. Int. ACM SIGIR Conf. Res. Develop.
Inf. Retr. Asia–Pacific Region, Nov. 2023, pp. 46–51.
[35] Y. Chen, Q. Fu, Y. Yuan, Z. Wen, G. Fan, D. Liu, D. Zhang, Z. Li, and
Y. Xiao, ‘‘Hallucination detection: Robustly discerning reliable answers
in large language models,’’ in Proc. 32nd ACM Int. Conf. Inf. Knowl.
Manage., Oct. 2023, pp. 245–255.
[36] T. Dutta and X. Liu, ‘‘FaCTQA: Detecting and localizing factual
errors in generated summaries through question and answering from
heterogeneous models,’’ in Proc. Int. Joint Conf. Neural Netw. (IJCNN),
Jun. 2024, pp. 1–8.
[37] K. Furumai, Y. Wang, M. Shinohara, K. Ikeda, Y. Yu, and T. Kato,
‘‘Detecting dialogue hallucination using graph neural networks,’’ in Proc.
Int. Conf. Mach. Learn. Appl. (ICMLA), Dec. 2023, pp. 871–877.
[38] R. Hu, J. Zhong, M. Ding, Z. Ma, and M. Chen, ‘‘Evaluation of
hallucination and robustness for large language models,’’ in Proc. IEEE
23rd Int. Conf. Softw. Qual., Rel., Secur. Companion (QRS-C), Oct. 2023,
pp. 374–382.
[39] S. Huo, N. Arabzadeh, and C. Clarke, ‘‘Retrieving supporting evi-
dence for generative question answering,’’ in Proc. Annu. Int. ACM
SIGIR Conf. Res. Develop. Inf. Retr. Asia–Pacific Region, Nov. 2023,
pp. 11–20.
[40] R. Koç, M. K. Gürkan, and F. T. Yarman Vural, ‘‘ReRag: A new
architecture for reducing the hallucination by retrieval- augmented
generation,’’ in Proc. 9th Int. Conf. Comput. Sci. Eng. (UBMK),
Oct. 2024, pp. 961–965.
[41] F. Li, X. Zhang, and P. Zhang, ‘‘Mitigating hallucination issues in small-
parameter LLMs through inter-layer contrastive decoding,’’ in Proc. Int.
Joint Conf. Neural Netw. (IJCNN), Jun. 2024, pp. 1–8.
[42] D. Li and F. Xu, ‘‘The deep integration of knowledge graphs and large
language models: Advancements, challenges, and future directions,’’ in
Proc. IEEE 2nd Int. Conf. Sensors, Electron. Comput. Eng. (ICSECE),
Aug. 2024, pp. 157–162.
[43] N. Ngu, N. Lee, and P. Shakarian, ‘‘Diversity measures: Domain-
independent proxies for failure in language model queries,’’ in Proc. IEEE
18th Int. Conf. Semantic Comput. (ICSC), Feb. 2024, pp. 176–182.
[44] P. Pezeshkpour, ‘‘Measuring and modifying factual knowledge in large
language models,’’ in Proc. Int. Conf. Mach. Learn. Appl. (ICMLA),
Dec. 2023, pp. 831–838.
[45] B. Snyder, M. Moisescu, and M. B. Zafar, ‘‘On early detection of hal-
lucinations in factual question answering,’’ in Proc. 30th ACM SIGKDD
Conf. Knowl. Discovery Data Mining, Aug. 2024, pp. 2721–2732.
[46] Z. Wei, D. Guo, D. Huang, Q. Zhang, S. Zhang, K. Jiang, and R.
Li, ‘‘Detecting and mitigating the ungrounded hallucinations in text
generation by LLMs,’’ in Proc. Int. Conf. Artif. Intell., Syst. Netw. Secur.,
Dec. 2023, pp. 77–81.
[47] X. Xia and S. Dong, ‘‘Optimizing inference capabilities in Chinese
NLP: A study on lightweight generative language models for knowledge
question answering,’’ in Proc. 4th Int. Conf. Neural Netw., Inf. Commun.
(NNICE), Jan. 2024, pp. 359–363.
[48] M. Alshammary, M. N. Uddin, and L. Khan, ‘‘RFPG: Question-
answering from low-resource language (Arabic) texts using factually
aware RAG,’’ in Proc. IEEE 10th Int. Conf. Collaboration Internet
Comput. (CIC), Oct. 2024, pp. 107–116.
[49] D. Chen, S. Wang, Z. Fan, X. Hu, and C. Li, ‘‘Freeze-CD: Alleviating
hallucination of large language models via contrastive decoding with
local freezing training,’’ in Proc. IEEE Int. Conf. Smart Internet Things
(SmartIoT), Nov. 2024, pp. 325–329.
[50] J. Feng, Q. Wang, H. Qiu, and L. Liu, ‘‘Retrieval in decoder benefits
generative models for explainable complex question answering,’’ Neural
Netw., vol. 181, pp. 1–14, Jan. 2025.
[51] C. N. Hang, P.-D. Yu, and C. W. Tan, ‘‘TrumorGPT: Query optimization
and semantic reasoning over networks for automated fact-checking,’’ in
Proc. 58th Annu. Conf. Inf. Sci. Syst. (CISS), Mar. 2024, pp. 1–6.
[52] K. Hu, M. Yan, W. H. Chong, Y. Keong Yap, C. Guan, J. T. Zhou, and
I. W. Tsang, ‘‘Ladder-of-thought: Using knowledge as steps to elevate
stance detection,’’ in Proc. Int. Joint Conf. Neural Netw. (IJCNN),
Jun. 2024, pp. 1–8.
[53] Q. Huang, F. Huang, D. Tao, Y. Zhao, B. Wang, and Y. Huang, ‘‘CoQ: AN
empirical framework for multi-hop question answering empowered by
large language models,’’ in Proc. IEEE Int. Conf. Acoust., Speech Signal
Process. (ICASSP), Apr. 2024, pp. 11566–11570.
[54] Y. Huang and G. Zeng, ‘‘RD-P: A trustworthy retrieval-augmented
prompter with knowledge graphs for LLMs,’’ in Proc. 33rd ACM Int.
Conf. Inf. Knowl. Manage., Oct. 2024, pp. 942–952.
[55] F. Shiri, V. Nguyen, F. Moghimifar, J. Yoo, G. Haffari, and Y.-F. Li,
‘‘Decompose, enrich, and extract! Schema-aware event extraction using
LLMs,’’ 2024, arXiv:2406.01045.
[56] Y. Song, H. Fan, J. Liu, Y. Liu, X. Ye, and Y. Ouyang, ‘‘A goal-
oriented document-grounded dialogue based on evidence generation,’’
Data Knowl. Eng., vol. 155, pp. 1–16, Jan. 2025.
[57] W. Su, Y. Tang, Q. Ai, C. Wang, Z. Wu, and Y. Liu, ‘‘Mitigating entity-
level hallucination in large language models,’’ in Proc. Annu. Int. ACM
SIGIR Conf. Res. Develop. Inf. Retr. Asia–Pacific Region, Dec. 2024,
pp. 23–31.
[58] D. Yang, J. Rao, K. Chen, X. Guo, Y. Zhang, J. Yang, and Y. Zhang,
‘‘IM-RAG: Multi-round retrieval-augmented generation through learning
inner monologues,’’ in Proc. 47th Int. ACM SIGIR Conf. Res. Develop.
Inf. Retr., Jul. 2024, pp. 730–740.
[59] H. Zhang, R. Zhang, J. Guo, M. de Rijke, Y. Fan, and X. Cheng, ‘‘Are
large language models good at utility judgments?’’ in Proc. 47th Int. ACM
SIGIR Conf. Res. Develop. Inf. Retr., Jul. 2024, pp. 1941–1951.
[60] Z. Zhang, J. Chen, W. Shi, L. Yi, C. Wang, and Q. Yu, ‘‘Contrastive
learning for knowledge-based question generation in large language
models,’’ in Proc. 5th Int. Conf. Intell. Comput. Hum.-Comput. Interact.
(ICHCI), Sep. 2024, pp. 583–587.
[61] Z. Zhu, G. Qi, G. Shang, Q. He, W. Zhang, N. Li, Y. Chen, L. Hu,
W. Zhang, and F. Dang, ‘‘Enhancing large language models with
knowledge graphs for robust question answering,’’ in Proc. IEEE 30th
Int. Conf. Parallel Distrib. Syst. (ICPADS), Oct. 2024, pp. 262–269.
[62] X. Shi, J. Liu, Y. Liu, Q. Cheng, and W. Lu, ‘‘Know where to go: Make
LLM a relevant, responsible, and trustworthy searchers,’’ Decis. Support
Syst., vol. 188, pp. 1–13, Jan. 2025.
[63] J. Sun, Y. Pan, and X. Yan, ‘‘Improving intermediate reasoning in zero-
shot chain-of-thought for large language models with filter supervisor-
self correction,’’ Neurocomputing, vol. 620, pp. 1–17, Mar. 2025.
[64] J. Wang, D. Cao, S. Lu, Z. Ma, J. Xiao, and T.-S. Chua, ‘‘Causal-driven
large language models with faithful reasoning for knowledge question
answering,’’ in Proc. 32nd ACM Int. Conf. Multimedia, Oct. 2024,
pp. 4331–4340.
[65] A. Kraft and E. Soulier, ‘‘Knowledge-enhanced language models are
not bias-proof: Situated knowledge and epistemic injustice in AI,’’ in
Proc. ACM Conf. Fairness, Accountability, Transparency, Jun. 2024,
pp. 1433–1445.
[66] M. Bahrami, R. Sonoda, and R. Srinivasan, ‘‘LLM diagnostic toolkit:
Evaluating LLMs for ethical issues,’’ in Proc. Int. Joint Conf. Neural
Netw. (IJCNN), Jun. 2024, pp. 1–8.
148250 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
[67] M. Ding, Y. Shen, and M. Chen, ‘‘Automated functionality and security
evaluation of large language models,’’ in Proc. 9th IEEE Int. Conf. Smart
Cloud (SmartCloud), May 2024, pp. 37–41.
[68] T. R. Hannigan, I. P. McCarthy, and A. Spicer, ‘‘Beware of botshit: How
to manage the epistemic risks of generative chatbots,’’ Bus. Horizons,
vol. 67, no. 5, pp. 471–486, Sep. 2024.
[69] Y. Huang, J. Song, Z. Wang, S. Zhao, H. Chen, F. Juefei-Xu, and L. Ma,
‘‘Look before you leap: An exploratory study of uncertainty analysis for
large language models,’’ IEEE Trans. Softw. Eng., vol. 51, no. 2, pp. 1–18,
Feb. 2025.
[70] K. Jiang, Q. Zhang, D. Guo, D. Huang, S. Zhang, Z. Wei, F. Ning, and
R. Li, ‘‘AI-generated news articles based on large language models,’’ in
Proc. Int. Conf. Artif. Intell., Syst. Netw. Secur., Dec. 2023, pp. 82–87.
[71] N. Li, Y. Li, Y. Liu, L. Shi, K. Wang, and H. Wang, ‘‘Drowzee:
Metamorphic testing for fact-conflicting hallucination detection in large
language models,’’ Proc. ACM Program. Lang., vol. 8, no. OOPSLA2,
pp. 1843–1872, Oct. 2024.
[72] S. B. Shah, S. Thapa, A. Acharya, K. Rauniyar, S. Poudel,
S. Jain, A. Masood, and U. Naseem, ‘‘Navigating the Web of
disinformation and misinformation: Large language models as
double-edged swords,’’ IEEE Access, early access, May 29, 2024,
doi: 10.1109/ACCESS.2024.3406644.
[73] S. Tripathi, H. Griffith, and H. Rathore, ‘‘Assessing hallucination in large
language models under adversarial attacks,’’ in Proc. 9th Int. Conf. Mobile
Secure Services (MobiSecServ), Nov. 2024, pp. 1–6.
[74] L. Wang, H. Zhang, H. Shao, M. Wu, and W. Ren, ‘‘WHW: An efficient
data organization method for fine-tuning large language models,’’ in
Proc. 5th Int. Conf. Inf. Sci., Parallel Distrib. Syst. (ISPDS), May 2024,
pp. 221–224.
[75] F. Perrina, F. Marchiori, M. Conti, and N. V. Verde, ‘‘AGIR: Automating
cyber threat intelligence reporting with natural language generation,’’ in
Proc. IEEE Int. Conf. Big Data (BigData), Dec. 2023, pp. 3053–3062.
[76] T. Yu, Y. Yao, H. Zhang, T. He, Y. Han, G. Cui, J. Hu, Z. Liu, H.-T. Zheng,
and M. Sun, ‘‘RLHF-V: Towards trustworthy MLLMs via behavior
alignment from fine-grained correctional human feedback,’’ in Proc.
IEEE/CVF Conf. Comput. Vis. Pattern Recognit. (CVPR), Jun. 2024,
pp. 13807–13816.
[77] C. Jiang, H. Xu, M. Dong, J. Chen, W. Ye, M. Yan, Q. Ye, J.
Zhang, F. Huang, and S. Zhang, ‘‘Hallucination augmented contrastive
learning for multimodal large language model,’’ in Proc. IEEE/CVF Conf.
Comput. Vis. Pattern Recognit. (CVPR), Jun. 2024, pp. 27026–27036.
[78] P. Kaul, Z. Li, H. Yang, Y. Dukler, A. Swaminathan, C. J. Taylor, and
S. Soatto, ‘‘THRONE: An object-based hallucination benchmark for
the free-form generations of large vision-language models,’’ in Proc.
IEEE/CVF Conf. Comput. Vis. Pattern Recognit. (CVPR), Jun. 2024,
pp. 27218–27228.
[79] A. Bendeck and J. Stasko, ‘‘An empirical evaluation of the GPT-4
multimodal language model on visualization literacy tasks,’’ IEEE Trans.
Vis. Comput. Graphics, vol. 31, no. 1, pp. 1105–1115, Jan. 2025.
[80] K. Bönisch, M. Stoeckel, and A. Mehler, ‘‘HyperCausal: Visualizing
causal inference in 3D hypertext,’’ in Proc. 35th ACM Conf. Hypertext
Social Media, Sep. 2024, pp. 330–336.
[81] P. Ding, J. Wu, J. Kuang, D. Ma, X. Cao, X. Cai, S. Chen, J. Chen,
and S. Huang, ‘‘Hallu-PI: Evaluating hallucination in multi-modal large
language models within perturbed inputs,’’ in Proc. 32nd ACM Int. Conf.
Multimedia, Oct. 2024, pp. 10707–10715.
[82] B. A. Halperin and S. M. Lukin, ‘‘Artificial dreams: Surreal visual
storytelling as inquiry into AI ‘Hallucination,’’’ in Proc. Designing
Interact. Syst. Conf., Jul. 2024, pp. 619–637.
[83] O. H. Hamid, ‘‘Beyond probabilities: Unveiling the delicate dance of large
language models (LLMs) and AI-hallucination,’’ in Proc. IEEE Conf.
Cognit. Comput. Aspects Situation Manage. (CogSIMA), May 2024,
pp. 85–90.
[84] C. Jiang, H. Jia, M. Dong, W. Ye, H. Xu, M. Yan, J. Zhang,
and S. Zhang, ‘‘Hal-Eval: A universal and fine-grained hallucination
evaluation framework for large vision language models,’’ in Proc. 32nd
ACM Int. Conf. Multimedia, Oct. 2024, pp. 525–534.
[85] S. Selva Kumar, A. K. M. A. Khan, I. A. Banday, M. Gada,
and V. V. Shanbhag, ‘‘Overcoming LLM challenges using RAG-driven
precision in coffee leaf disease remediation,’’ in Proc. Int. Conf. Emerg.
Technol. Comput. Sci. Interdiscipl. Appl. (ICETCS), Apr. 2024, pp. 1–6.
[86] Y. Li, X. Lan, H. Chen, K. Lu, and D. Jiang, ‘‘Multimodal PEAR chain-
of-thought reasoning for multimodal sentiment analysis,’’ ACM Trans.
Multimedia Comput., Commun., Appl., vol. 20, no. 9, pp. 1–23, Sep. 2024.
[87] F. Ma, X. Jin, H. Wang, Y. Xian, J. Feng, and Y. Yang, ‘‘Vista-llama:
Reducing hallucination in video language models via equal distance to
visual tokens,’’ in Proc. IEEE/CVF Conf. Comput. Vis. Pattern Recognit.
(CVPR), Jun. 2024, pp. 13151–13160.
[88] A. A. Verma, A. Saeidi, S. Hegde, A. Therala, F. D. Bardoliya,
N. Machavarapu, S. A. K. Ravindhiran, S. Malyala, A. Chatterjee, Y.
Yang, and C. Baral, ‘‘Evaluating multimodal large language models
across distribution shifts and augmentations,’’ in Proc. IEEE/CVF
Conf. Comput. Vis. Pattern Recognit. Workshops (CVPRW), Jun. 2024,
pp. 5314–5324.
[89] Y. Wang, T. Wang, and L. Zhang, ‘‘Consistency framework for zero-shot
image captioning,’’ in Proc. 4th Int. Conf. Neural Netw., Inf. Commun.
(NNICE), Jan. 2024, pp. 511–515.
[90] Y. Qian, J. Wei, Y. Zhang, X. Zhang, C. Wei, S. Chen, Y. Li,
C. Ye, B. Huang, and H. Wang, ‘‘CGSMP: Controllable generative
summarization via multimodal prompt,’’ in Proc. 1st Workshop Large
Generative Models Meet Multimodal Appl., Nov. 2023, pp. 45–50.
[91] J. Yu, Y. Zhang, Z. Zhang, Z. Yang, G. Zhao, F. Sun, F. Zhang, Q. Liu,
J. Sun, J. Liang, and Y. Zhang, ‘‘RAG-guided large language models for
visual spatial description with adaptive hallucination corrector,’’ in Proc.
32nd ACM Int. Conf. Multimedia, Oct. 2024, pp. 11407–11413.
[92] J. Fei, T. Wang, J. Zhang, Z. He, C. Wang, and F. Zheng, ‘‘Transferable
decoding with visual entities for zero-shot image captioning,’’ in Proc.
IEEE/CVF Int. Conf. Comput. Vis. (ICCV), Oct. 2023, pp. 3113–3123.
[93] H. Fei, M. Luo, J. Xu, S. Wu, W. Ji, M.-L. Lee, and W. Hsu, ‘‘Fine-grained
structural hallucination detection for unified visual comprehension and
generation in multimodal LLM,’’ in Proc. 1st ACM Multimedia Workshop
Multi-modal Misinformation Governance Era Found. Models, Oct. 2024,
pp. 13–22.
[94] T. Gao, P. Chen, M. Zhang, C. Fu, Y. Shen, Y. Zhang, S. Zhang, X. Zheng,
X. Sun, L. Cao, and R. Ji, ‘‘Cantor: Inspiring multimodal chain-of-
thought of MLLM,’’ in Proc. 32nd ACM Int. Conf. Multimedia, Oct. 2024,
pp. 9096–9105.
[95] Q. Huang, X. Dong, P. Zhang, B. Wang, C. He, J. Wang, D. Lin, W. Zhang,
and N. Yu, ‘‘OPERA: Alleviating hallucination in multi-modal large
language models via over-trust penalty and retrospection-allocation,’’
in Proc. IEEE/CVF Conf. Comput. Vis. Pattern Recognit. (CVPR),
Jun. 2024, pp. 13418–13427.
[96] K. Ranasinghe, S. N. Shukla, O. Poursaeed, M. S. Ryoo, and T.-Y. Lin,
‘‘Learning to localize objects improves spatial reasoning in visual-
LLMs,’’ in Proc. IEEE/CVF Conf. Comput. Vis. Pattern Recognit.
(CVPR), Jun. 2024, pp. 12977–12987.
[97] N. Sarfati, I. Yerushalmy, M. Chertok, and Y. Keller, ‘‘Generating
factually consistent sport highlights narrations,’’ in Proc. 6th Int.
Workshop Multimedia Content Anal. Sports, Oct. 2023, pp. 15–22.
[98] Y. Wang, Y. Zeng, J. Liang, X. Xing, J. Xu, and X. Xu, ‘‘RetrievalMMT:
Retrieval-constrained multi-modal prompt learning for multi-modal
machine translation,’’ in Proc. Int. Conf. Multimedia Retr., May 2024,
pp. 860–868.
[99] C. Wang, J. Yang, Y. Zhou, and X. Yue, ‘‘CooKie: Commonsense
knowledge-guided mixture-of-experts framework for fine-
grained visual question answering,’’ Inf. Sci., vol. 695, pp. 1–20,
Mar. 2025.
[100] J. Wei and X. Zhang, ‘‘DOPRA: Decoding over-accumulation penaliza-
tion and re-allocation in specific weighting layer,’’ in Proc. 32nd ACM
Int. Conf. Multimedia, Oct. 2024, pp. 7065–7074.
[101] Q. Yu, J. Li, L. Wei, L. Pang, W. Ye, B. Qin, S. Tang, Q. Tian,
and Y. Zhuang, ‘‘HalluciDoctor: Mitigating hallucinatory toxicity in
visual instruction data,’’ in Proc. IEEE/CVF Conf. Comput. Vis. Pattern
Recognit. (CVPR), Jun. 2024, pp. 12944–12953.
[102] M. Gao, S. Chen, L. Pang, Y. Yao, J. Dang, W. Zhang, J. Li, S. Tang,
Y. Zhuang, and T. Chua, ‘‘Fact: Teaching MLLMs with faithful, concise
and transferable rationales,’’ in Proc. 32nd ACM Int. Conf. Multimedia,
Oct. 2024, pp. 846–855.
[103] P. Xu, W. Shao, K. Zhang, P. Gao, S. Liu, M. Lei, F. Meng,
S. Huang, Y. Qiao, and P. Luo, ‘‘LVLM-EHub: A comprehen-
sive evaluation benchmark for large vision-language models,’’ IEEE
Trans. Pattern Anal. Mach. Intell., vol. 47, no. 3, pp. 1877–1893,
Mar. 2025.
VOLUME 13, 2025 148251
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
[104] Y. Zhang, Z. Ma, X. Gao, S. Shakiah, Q. Gao, and J. Chai, ‘‘Groundhog
grounding large language models to holistic segmentation,’’ in Proc.
IEEE/CVF Conf. Comput. Vis. Pattern Recognit. (CVPR), Jun. 2024,
pp. 14227–14238.
[105] J. Kang, W. Pan, T. Zhang, Z. Wang, S. Yang, Z. Wang, J. Wang, and
X. Niu, ‘‘Correcting factuality hallucination in complaint large language
model via entity-augmented,’’ in Proc. Int. Joint Conf. Neural Netw.
(IJCNN), Jun. 2024, pp. 1–8.
[106] C.-H.-J. Leung, Y. Yi, L. Kuai, Z. Li, S.-K.-A. Yeung, K.-W.-J. Lee,
K.-H.-K. Ho, and K. Hung, ‘‘RAG for question-answering for vocal
training based on domain knowledge base,’’ in Proc. 11th Int. Conf.
Behavioural Social Comput. (BESC), Aug. 2024, pp. 1–6.
[107] P. Ruangchutiphophan, C. Saetia, T. S. N. Ayutthaya, and T. Chalothorn,
‘‘Thai knowledge-augmented language model adaptation (ThaiKALA),’’
in Proc. 18th Int. Joint Symp. Artif. Intell. Natural Lang. Process. (iSAI-
NLP), Nov. 2023, pp. 1–6.
[108] L. Vidyaratne, X. Y. Lee, A. Kumar, T. Watanabe, A. Farahat, and
C. Gupta, ‘‘Generating troubleshooting trees for industrial equipment
using large language models (LLM),’’ in Proc. IEEE Int. Conf.
Prognostics Health Manage. (ICPHM), Jun. 2024, pp. 116–125.
[109] J. Yang, H. Xu, R. Wang, X. Ming, and S. Li, ‘‘Generating evaluation
criteria of domain-specific large language model using word vector
clustering,’’ in Proc. IEEE 24th Int. Conf. Softw. Qual., Rel., Secur.
Companion (QRS-C), Jul. 2024, pp. 94–100.
[110] S. Yao, Q. Ke, K. Li, Q. Wang, and J. Hu, ‘‘News GPT: A large
language model for reliable and hallucination-controlled news genera-
tion,’’ in Proc. 3rd Int. Symp. Robot., Artif. Intell. Inf. Eng., Jul. 2024,
pp. 113–119.
[111] Y. Bei, Z. Fang, S. Mao, S. Yu, Y. Jiang, Y. Tong, and W. Cai,
‘‘Manufacturing domain QA with integrated term enhanced RAG,’’ in
Proc. Int. Joint Conf. Neural Netw. (IJCNN), Jun. 2024, pp. 1–8.
[112] S. Roychowdhury, A. Alvarez, B. Moore, M. Krema, M. P. Gelpi,
P. Agrawal, F. M. Rodríguez, Á. Rodríguez, J. R. Cabrejas, P. M. Serrano,
and A. Mukherjee, ‘‘Hallucination-minimized data-to-answer framework
for financial decision-makers,’’ in Proc. IEEE Int. Conf. Big Data
(BigData), Dec. 2023, pp. 4693–4702.
[113] J. Sheng, ‘‘An augmentable domain-specific models for financial
analysis,’’ in Proc. 16th Int. Congr. Image Signal Process., Biomed. Eng.
Informat. (CISP-BMEI), Oct. 2023, pp. 1–4.
[114] K. Shetty, S. K. Bojanki, and A. Ratnaparkhi, ‘‘Sovereign risk summariza-
tion,’’ in Proc. 5th ACM Int. Conf. AI Finance, Nov. 2024, pp. 779–786.
[115] S. Wang, B. Xie, L. Ding, X. Gao, J. Chen, and Y. Xiang, ‘‘SeCor:
Aligning semantic and collaborative representations by large language
models for next-point-of-interest recommendations,’’ in Proc. 18th ACM
Conf. Recommender Syst., Oct. 2024, pp. 1–11.
[116] L. Xu and J. Liu, ‘‘A chat bot for enrollment of Xi’an jiaotong-liverpool
university based on RAG,’’ in Proc. 8th Int. Workshop Control Eng. Adv.
Algorithms (IWCEAA), Nov. 2024, pp. 125–129.
[117] J. Zheng, H. Wang, and J. Yao, ‘‘Building lightweight domain-
specific consultation systems via inter-external knowledge fusion
contrastive learning,’’ IEEE Access, vol. 12, pp. 113244–113258,
2024.
[118] J. Zheng, F. Xu, W. Chen, Z. Fang, and J. Yao, ‘‘Core-view contrastive
learning network for building lightweight cross-domain consultation
system,’’ IEEE Access, vol. 12, pp. 65615–65629, 2024.
[119] B. Zhan, Y. Duan, and S. Yan, ‘‘IC-BERT: An instruction classifier
model alleviates the hallucination of large language models in traditional
Chinese medicine,’’ in Proc. IEEE 9th Int. Conf. Comput. Intell. Appl.
(ICCIA), Aug. 2024, pp. 221–225.
[120] J. J. Chang and E. Y. Chang, ‘‘SocraHealth: Enhancing medical diagnosis
and correcting historical records,’’ in Proc. Int. Conf. Comput. Sci.
Comput. Intell. (CSCI), Dec. 2023, pp. 1400–1405.
[121] D. B. Craig and S. Drăghici, ‘‘What’s the data say? An LLM-based system
for interrogating experimental data,’’ in Proc. IEEE Int. Conf. Bioinf.
Biomed. (BIBM), Dec. 2024, pp. 1457–1462.
[122] E. Herron, J. Yin, and F. Wang, ‘‘SciTrust: Evaluating the trustworthiness
of large language models for science,’’ in Proc. SC24-W: Workshops
Int. Conf. High Perform. Comput., Netw., Storage Anal., Nov. 2024,
pp. 72–78.
[123] T. Pang, K. Tan, Y. Yao, X. Liu, F. Meng, C. Fan, and X. Zhang,
‘‘REMED: Retrieval-augmented medical document query responding
with embedding fine-tuning,’’ in Proc. Int. Joint Conf. Neural Netw.
(IJCNN), Jun. 2024, pp. 1–8.
[124] A. Tariq, S. Trivedi, A. U. Khan, G. Ramasamy, S. Fathizadeh, M. T. Stib,
N. Tan, B. N. Patel, and I. Banerjee, ‘‘Patient-centric summarization of
radiology findings using two-step training of large language models,’’ in
Proc. ACM Trans. Comput. Healthcare, 2024, pp. 1–16.
[125] Z. Zeng, Q. Cheng, X. Hu, Y. Zhuang, X. Liu, K. He, and Z. Liu,
‘‘KoSEL: Knowledge subgraph enhanced large language model for
medical question answering,’’ Knowl.-Based Syst., vol. 309, pp. 1–18,
Jan. 2025.
[126] J. Zhao, Q. Guo, J. Liang, Z. Li, and Y. Xiao, ‘‘Effective in-context
learning for named entity recognition,’’ in Proc. IEEE Int. Conf. Bioinf.
Biomed. (BIBM), Dec. 2024, pp. 1376–1382.
[127] V. Sharma and V. Raman, ‘‘A reliable knowledge processing framework
for combustion science using foundation models,’’ Energy AI, vol. 16,
pp. 1–21, May 2024.
[128] D. Xu, Z. Zhang, Z. Zhu, Z. Lin, Q. Liu, X. Wu, T. Xu, W. Wang,
Y. Ye, X. Zhao, E. Chen, and Y. Zheng, ‘‘Editing factual knowledge and
explanatory ability of medical large language models,’’ in Proc. 33rd
ACM Int. Conf. Inf. Knowl. Manage., Oct. 2024, pp. 2660–2670.
[129] F. Leiser, S. Eckhardt, M. Knaeble, A. Maedche, G. Schwabe, and
A. Sunyaev, ‘‘From ChatGPT to FactGPT: A participatory design study
to mitigate the effects of large language model hallucinations on users,’’
in Proc. Mensch und Comput., Sep. 2023, pp. 81–90.
[130] F. Leiser, S. Eckhardt, V. Leuthe, M. Knaeble, A. Mädche, G. Schwabe,
and A. Sunyaev, ‘‘HILL: A hallucination identifier for large language
models,’’ in Proc. CHI Conf. Human Factors Comput. Syst., May 2024,
pp. 1–13.
[131] J. Li, R. Yuan, Y. Tian, and J. Li, ‘‘Towards instruction-tuned verification
for improving biomedical information extraction with large language
models,’’ in Proc. IEEE Int. Conf. Bioinf. Biomed. (BIBM), Dec. 2024,
pp. 6685–6692.
[132] R. Y. Maragheh, C. Fang, C. C. Irugu, P. Parikh, J. Cho, J. Xu, S. Sukumar,
M. Patel, E. Körpeoglu, S. Kumar, and K. Achan, ‘‘LLM-TAKE: Theme-
aware keyword extraction using large language models,’’ in Proc. IEEE
Int. Conf. Big Data (BigData), Apr. 2023, pp. 4318–4324.
[133] T. T. Procko, A. Davidoff, T. Elvira, and O. Ochoa, ‘‘Leveraging
large language models on the traditional scientific writing work-
flow,’’ in Proc. Conf. AI, Sci., Eng., Technol. (AIxSET), Sep. 2024,
pp. 154–161.
[134] J. Waldo and S. Boussard, ‘‘GPTs and hallucination,’’ Commun. ACM,
vol. 68, no. 1, pp. 40–45, 2024.
[135] S. Zhao and X. Sun, ‘‘Enabling controllable table-to-text generation via
prompting large language models with guided planning,’’ Knowledge-
Based Syst., vol. 304, pp. 1–9, Nov. 2024.
[136] F. Cheng, V. Zouhar, S. Arora, M. Sachan, H. Strobelt, and M. El-Assady,
‘‘RELIC: Investigating large language model responses using self-
consistency,’’ in Proc. CHI Conf. Human Factors Comput. Syst.,
May 2024, pp. 1–18.
[137] E. Lin, J. Hale, and J. Gratch, ‘‘Toward a better understand-
ing of the emotional dynamics of negotiation with large language
models,’’ in Proc. Twenty-fourth Int. Symp. Theory, Algorithmic
Found., Protocol Design Mobile Netw. Mobile Comput., Oct. 2023,
pp. 545–550.
[138] Y. Zhao, J. Wu, X. Wang, W. Tang, D. Wang, and M. de Rijke, ‘‘Let me do
it for you: Towards LLM empowered recommendation via tool learning,’’
in Proc. 47th Int. ACM SIGIR Conf. Res. Develop. Inf. Retr., Jul. 2024,
pp. 1796–1806.
[139] Z. Sun, K. Feng, J. Yang, X. Qu, H. Fang, Y.-S. Ong, and W. Liu,
‘‘Adaptive in-context learning with large language models for bundle
generation,’’ in Proc. 47th Int. ACM SIGIR Conf. Res. Develop. Inf. Retr.,
Jul. 2024, pp. 966–976.
[140] X. Zhang, R. Xie, Y. Lyu, X. Xin, P. Ren, M. Liang, B. Zhang, Z. Kang, M.
de Rijke, and Z. Ren, ‘‘Towards empathetic conversational recommender
systems,’’ in Proc. 18th ACM Conf. Recommender Syst., Oct. 2024,
pp. 84–93.
[141] Y. Li, Y. He, R. Lian, and Q. Guo, ‘‘Fault diagnosis and system
maintenance based on large language models and knowledge graphs,’’
in Proc. 5th Int. Conf. Robot., Intell. Control Artif. Intell. (RICAI),
Dec. 2023, pp. 589–592.
[142] K. Li, Q. Zheng, Y. Zhan, C. Zhang, T. Zhang, X. Lin, C. Qi, L. Li,
and D. Tao, ‘‘Alleviating action hallucination for LLM-based embodied
agents via inner and outer alignment,’’ in Proc. 7th Int. Conf. Pattern
Recognit. Artif. Intell. (PRAI), Aug. 2024, pp. 613–621.
148252 VOLUME 13, 2025
C. Woesle et al.: Systematic Literature Review of Hallucinations in Large Language Models
[143] K. Manas, S. Zwicklbauer, and A. Paschke, ‘‘CoT-TL: Low-resource
temporal knowledge representation of planning instructions using chain-
of-thought reasoning,’’ in Proc. IEEE/RSJ Int. Conf. Intell. Robots Syst.
(IROS), Oct. 2024, pp. 9636–9643.
[144] J. Chen and S. Lu, ‘‘An advanced driving agent with the multimodal
large language model for autonomous vehicles,’’ in Proc. IEEE Int. Conf.
Mobility, Oper., Services Technol. (MOST), May 2024, pp. 1–11.
[145] S. Jha, S. K. Jha, P. Lincoln, N. D. Bastian, A. Velasquez, and S. Neema,
‘‘Dehallucinating large language models using formal methods guided
iterative prompting,’’ in Proc. IEEE Int. Conf. Assured Autonomy (ICAA),
Jun. 2023, pp. 149–152.
[146] H. Nakajima and J. Miura, ‘‘Combining ontological knowledge and large
language model for user-friendly service robots,’’ in Proc. IEEE/RSJ Int.
Conf. Intell. Robots Syst. (IROS), Oct. 2024, pp. 4755–4762.
[147] Y. Qiu, ‘‘The impact of LLM hallucinations on motor skill
learning: A case study in badminton,’’ IEEE Access, vol. 12,
pp. 139669–139682, 2024.
[148] W. de Almeida da Silva, L. C. Costa Fonseca, S. Labidi, and J. C. Lima
Pacheco, ‘‘Mitigation of hallucinations in language models in education:
A new approach of comparative and cross-verification,’’ in Proc. IEEE
Int. Conf. Adv. Learn. Technol. (ICALT), Jul. 2024, pp. 207–209.
[149] H.-T. Ho, D.-T. Ly, and L. V. Nguyen, ‘‘Mitigating hallucinations in large
language models for educational application,’’ in Proc. IEEE Int. Conf.
Consum. Electronics-Asia (ICCE-Asia), Nov. 2024, pp. 1–4.
[150] G. Wang, S. Qin, and X. Liu, ‘‘Exploration and practice of applying large
language models in home education guidance,’’ in Proc. 4th Int. Conf. Inf.
Sci. Educ. (ICISE-IE), Dec. 2023, pp. 201–208.
[151] J. Shi, J. Zhao, X. Wu, R. Xu, Y.-H. Jiang, and L. He, ‘‘Mitigating
reasoning hallucination through multi-agent collaborative filtering,’’
Expert Syst. Appl., vol. 263, Mar. 2025, Art. no. 125723.
[152] N. Dziri, E. Kamalloo, S. Milton, O. Zaiane, M. Yu, E. M. Ponti, and
S. Reddy, ‘‘FaithDial: A faithful benchmark for information-seeking
dialogue,’’ Trans. Assoc. for Comput. Linguistics, vol. 10, pp. 1473–1490,
Dec. 2022.
[153] J. Lee, S. Ahn, D. Kim, and D. Kim, ‘‘Performance comparison of
retrieval-augmented generation and fine-tuned large language models
for construction safety management knowledge retrieval,’’ Autom.
Construct., vol. 168, pp. 1–12, Dec. 2024.
[154] D. Kahneman, Thinking, Fast and Slow. Baltimore, MD, USA: Penguin,
2012.
[155] J. S. B. T. Evans, ‘‘Dual-processing accounts of reasoning, judgment,
and social cognition,’’ Annu. Rev. Psychol., vol. 59, no. 1, pp. 255–278,
Jan. 2008.
[156] K. E. Stanovich and R. F. West, ‘‘Individual differences in reasoning:
Implications for the rationality debate?’’ Behav. Brain Sci., vol. 23, no. 5,
pp. 645–665, Oct. 2000.
[157] G. Norman, J. Sherbino, K. Dore, T. Wood, M. Young, W. Gaissmaier,
S. Kreuger, and S. Monteiro, ‘‘The etiology of diagnostic errors:
A controlled trial of system 1 versus system 2 reasoning,’’ Academic
Med., vol. 89, no. 2, pp. 277–284, Feb. 2014.
[158] X. Cheng, J. Li, W. Xin Zhao, and J.-R. Wen, ‘‘Think more, hallucinate
less: Mitigating hallucinations via dual process of fast and slow thinking,’’
2025, arXiv:2501.01306.
[159] S. C. Bellini-Leite, ‘‘Dual process theory for large language models: An
overview of using psychology to address hallucination and reliability
issues,’’ Adapt. Behav., vol. 32, no. 4, pp. 329–343, Aug. 2024.
CHRISTIAN WOESLE received the B.Sc. degree
in business informatics from the FOM University
of Applied Sciences, Munich, Germany, and the
master’s degree in systems engineering from the
University of Applied Sciences, Munich. He is
currently pursuing the Ph.D. degree with the
Chair of Hybrid Intelligence, Helmut-Schmidt-
University/University of the Federal Armed Forces
Hamburg. His current research interests include
hallucinations in LLMs, multi-agent systems, and
natural language processing.
LEOPOLD FISCHER-BRANDIES (Member,
IEEE) received the B.Sc. degree in business
administration and the M.Sc. degree in digitaliza-
tion and entrepreneurship from the University of
Bayreuth, Germany. He is currently pursuing the
Ph.D. degree with the Chair of Hybrid Intelligence,
Helmut-Schmidt-University/University of the
Federal Armed Forces Hamburg. His current
research interests include deep learning, multi-
agent systems, natural language processing, and
detecting deepfakes.
RICARDO BUETTNER (Senior Member, IEEE)
received the Dipl.-Inf. degree in computer science
and the Dipl.-Wirtsch.-Ing. degree in industrial
engineering and management from the Technical
University of Ilmenau, Germany, the Dipl.-Kfm.
degree in business administration from the Uni-
versity of Hagen, Germany, the Ph.D. degree
in information systems from the University of
Hohenheim, Germany, and the Habilitation (venia
legendi) degree in information systems from the
University of Trier, Germany. He is currently a Chaired Professor of Hybrid
Intelligence with the Helmut-Schmidt-University/University of the Federal
Armed Forces Hamburg, Germany. He has published more than 160 peer-
reviewed articles, including articles in Electronic Markets, AIS Transactions
on Human-Computer Interaction, Personality and Individual Differences,
European Journal of Psychological Assessment, PLOS One, and IEEE
ACCESS. He has received 21 international best paper, best reviewer, and
service awards, and award nominations, including best paper awards by
AIS Transactions on Human-Computer Interaction, Electronic Markets, and
HICSS, for his work.
VOLUME 13, 2025 148253