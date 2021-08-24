# GoGirl Workshop Badge Generator

This document is to provide a step by step guide on how to generate a customised badge for the GoGirlGo4IT attendees


## Pre-requisites

 -  badge.html
 -  JPG image to use in the badge (publicly accessible)
 -  public hosting mechanism (e.g. S3/GCS)

## Steps to get your badge generator ready 

 -  **Step 1**: Edit badge.html
	 - Change the value of hidden input element "inFile" to the JPG image you want to use for your badge
	 - Make sure all other form element names are kept as is, including the on submit action

 -  **Step 2**: Copy your updated badge.html at a location where it can be accessed by the attendees (e.g. a publicly accessible S3 or GCS bucket)
 
 -  **Step 3**: You're all set! Just share the address of your badge.html with the attendees

## How it all works?

The form takes the name of the participant as input and submits it to our software running on AWS. This software puts the name at the bottom of the image. Check out our sample "gogirl-devops-badge.jpg for reference.

## Contact in case of questions:

 - bhavana.katara@anz.com
 - sirisha.vadrevu@anz.com
